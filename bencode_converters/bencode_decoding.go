package bencode_converters

import (
	"errors"
	"fmt"
	"strings"
)

func decodeNextElement(bytes []byte, index int) (string, int, error) {
	if bytes[index] == 'i' {
		return decodeNumber(bytes, index)
	} else if bytes[index] == 'd' {
		return decodeDict(bytes, index)
	} else if bytes[index] == 'l' {
		return decodeList(bytes, index)
	} else if byteIsDigit(bytes[index]) {
		return decodeStringLiteral(bytes, index)
	} else {
		return "", len(bytes), errors.New(fmt.Sprintf("No known entity start for index: %d, symbol: '%s'", index, string(bytes[index])))
	}
}

func decodeList(bytes []byte, index int) (string, int, error) {
	startIndex := index
	index++
	builder := strings.Builder{}
	builder.WriteString("[")
	var tempStr string
	var err error
	for index < len(bytes) && bytes[index] != 'e' {
		tempStr, index, err = decodeNextElement(bytes, index)
		if err != nil {
			return "", 0, err
		}
		builder.WriteString(tempStr)
		if bytes[index] != 'e' {
			builder.WriteString(",")
		}
	}
	if index >= len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("No closing 'e' symbol for list starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	builder.WriteString("]")
	return builder.String(), index + 1, err // +1 for 'e'
}

func decodeNumber(bytes []byte, index int) (string, int, error) {
	startIndex := index
	index++
	builder := strings.Builder{}
	for index < len(bytes) && bytes[index] != 'e' {
		builder.WriteByte(bytes[index])
		index++
	}
	if index >= len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("No closing 'e' symbol for number starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	return builder.String(), index + 1, nil // +1 for 'e'
}

func decodeStringLiteral(bytes []byte, index int) (string, int, error) {
	strLen := 0
	for byteIsDigit(bytes[index]) {
		strLen = strLen*10 + int(bytes[index]-'0')
		index++
	}
	if strLen <= 0 {
		return "", 0, errors.New(fmt.Sprintf("String literal len must be positive number. Index: %d, Symbol: '%s'", index, string(bytes[index])))
	} else if bytes[index] != ':' {
		return "", 0, errors.New(fmt.Sprintf("Missing string semicolon. Index: %d, Symbol: '%s'", index, string(bytes[index])))
	} else if index+strLen > len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("Wrong string length. Index: %d, Symbol: '%s'", index, string(bytes[index])))
	}
	index++ // semicolon
	resultingStr := "\"" + string(bytes[index:index+strLen]) + "\""
	return resultingStr, index + strLen, nil
}

func decodeDict(bytes []byte, index int) (string, int, error) {
	startIndex := index
	index++
	builder := strings.Builder{}
	builder.WriteString("{")
	readingKey := true
	var value, key string
	var err error
	for index < len(bytes) && bytes[index] != 'e' {
		if readingKey && byteIsDigit(bytes[index]) {
			key, index, err = decodeStringLiteral(bytes, index)
			if err != nil {
				return "", 0, err
			}
			builder.WriteString(key)
			builder.WriteString(":")
			readingKey = !readingKey
		} else if readingKey {
			return "", 0, errors.New(fmt.Sprintf("Dictionary key can be nothing but string. Index: %d, Symbol: '%s'", index, string(bytes[index])))
		} else if !readingKey {
			value, index, err = decodeNextElement(bytes, index)
			if err != nil {
				return "", 0, err
			}
			builder.WriteString(value)
			if bytes[index] != 'e' {
				builder.WriteString(",")
			}
			readingKey = !readingKey
		}
	}
	if index >= len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("No closing 'e' symbol for dictionary starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	builder.WriteString("}")
	return builder.String(), index + 1, err // +1 for 'e'
}

func BencodeToJSON(encodedStr []byte) (string, error) {
	var decodedStr string
	var err error
	if len(encodedStr) > 0 {
		decodedStr, _, err = decodeNextElement(encodedStr, 0)
		return decodedStr, err
	} else {
		return "", nil
	}
}
