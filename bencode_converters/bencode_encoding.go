package bencode_converters

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func byteIsDigit(symbol byte) bool {
	return 47 < symbol && symbol < 58 || symbol == 45
}

func encodeNextElement(json string, index int) ([]byte, int, error) {
	if byteIsDigit(json[index]) {
		return encodeNumber(json, index)
	} else if json[index] == '{' {
		return encodeDict(json, index)
	} else if json[index] == '[' {
		return encodeList(json, index)
	} else if json[index] == '"' {
		return encodeStringLiteral(json, index)
	}
	return nil, 0, errors.New(fmt.Sprintf("No known entity start for index: %d, symbol: '%s'", index, string(json[index])))
}

func encodeDict(json string, index int) ([]byte, int, error) {
	startIndex := index
	index++
	readingKey := true
	var tempBytes, resultingBytes []byte
	var err error
	var keyMap = map[string][]byte{}
	var currentKey string
	resultingBytes = append(resultingBytes, 'd')
	for index < len(json) && json[index] != '}' {
		if readingKey && json[index] == ':' {
			readingKey = !readingKey
		} else if readingKey && json[index] == '"' {
			var newIndex int
			tempBytes, newIndex, err = encodeStringLiteral(json, index)
			if err != nil {
				return nil, 0, err
			}
			currentKey = json[index+1:newIndex]
			keyMap[currentKey] = tempBytes
			index = newIndex
		} else if readingKey && json[index] != ',' {
			return nil, 0, errors.New(fmt.Sprintf("Dictionary key can be nothing but string. Index: %d, Symbol: '%s'", index, string(json[index])))
		} else if !readingKey {
			tempBytes, index, err = encodeNextElement(json, index)
			if err != nil {
				return nil, 0, err
			}

			keyMap[currentKey] = append(keyMap[currentKey], tempBytes...)
			readingKey = !readingKey
		}
		index++
	}
	if index >= len(json) {
		return nil, 0, errors.New(fmt.Sprintf("No closing '}' symbol for dictionary starting at index: %d", startIndex))
	}

	resultingBytes = append(resultingBytes, lexicographicConcat(keyMap)...)
	resultingBytes = append(resultingBytes, 'e')
	return resultingBytes, index, err
}

func lexicographicConcat(m map[string][]byte) []byte {
	var result []byte
	for len(m) > 0 {
		var minKey *string
		for key := range m {
			if minKey == nil || strings.Compare(key, *minKey) <= 0 {
				temp := key
				minKey = &temp
			}
		}
		result = append(result, m[*minKey]...)
		delete(m, *minKey)
	}
	return result
}

func encodeList(json string, index int) ([]byte, int, error) {
	startIndex := index
	index++
	var tempBytes, resultingBytes []byte
	var err error
	resultingBytes = append(resultingBytes, 'l')
	for index < len(json) && json[index] != ']' {
		if json[index] != ',' {
			tempBytes, index, err = encodeNextElement(json, index)
			resultingBytes = append(resultingBytes, tempBytes...)
			if err != nil {
				return nil, 0, err
			}
		}
		index++
	}
	if index >= len(json) {
		return nil, 0, errors.New(fmt.Sprintf("No closing ']' symbol for list starting at index: %d", startIndex))
	}
	resultingBytes = append(resultingBytes, 'e')
	return resultingBytes, index, err
}

func encodeStringLiteral(json string, index int) ([]byte, int, error) {
	index++
	literalEnd := strings.IndexByte(json[index:], '"')
	literal := json[index : index+literalEnd]
	var resultingBytes, temp []byte
	var err error
	var length int
	temp = []byte(literal)
	if err != nil {
		return nil, 0, err
	}
	length = len(temp)
	resultingBytes = append(resultingBytes, strconv.Itoa(length)...)
	resultingBytes = append(resultingBytes, ':')
	resultingBytes = append(resultingBytes, temp...)
	return resultingBytes, index + literalEnd, nil
}

func encodeNumber(json string, index int) ([]byte, int, error) {
	builder := strings.Builder{}
	builder.WriteByte('i')
	for index < len(json) && (byteIsDigit(json[index]) || json[index] == '-') {
		builder.WriteByte(json[index])
		index++
	}
	if index < len(json) && json[index] == '.' {
		return nil, 0, errors.New(fmt.Sprintf("Numbers must be integer. Index: %d", index))
	}
	builder.WriteByte('e')
	return []byte(builder.String()), index - 1, nil
}

func JSONToBencode(json string) ([]byte, error) {
	var encodedBytes []byte
	var err error
	if len(json) > 0 {
		encodedBytes, _, err = encodeNextElement(json, 0)
		return encodedBytes, err
	}
	return encodedBytes, err
}
