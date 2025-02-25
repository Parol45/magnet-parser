package bencode_obj

import (
	"errors"
	"fmt"
	"magnet-parser/globals"
	"strconv"
	"strings"
)

func decodeNumber(bytes []byte, index int) (int, int, error) {
	startIndex := index
	index++
	builder := strings.Builder{}
	for index < len(bytes) && bytes[index] != 'e' {
		builder.WriteByte(bytes[index])
		index++
	}
	if index >= len(bytes) {
		return 0, 0, errors.New(fmt.Sprintf("No closing 'e' symbol for number starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	res, err := strconv.Atoi(builder.String())
	if err != nil {
		return 0, 0, errors.New(fmt.Sprintf("Bad number format at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	return res, index + 1, nil // +1 for 'e'
}

func decodeStringList(bytes []byte, index int) ([]string, int, error) {
	startIndex := index
	index++
	var res []string
	var tempStr string
	var err error
	for index < len(bytes) && bytes[index] != 'e' {
		tempStr, index, err = decodeStringLiteral(bytes, index)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, tempStr)
	}
	if index >= len(bytes) {
		return nil, 0, errors.New(fmt.Sprintf("No closing 'e' symbol for list starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	return res, index + 1, err // +1 for 'e'
}

func decodeStringLiteral(bytes []byte, index int) (string, int, error) {
	strLen := 0
	for globals.ByteIsDigit(bytes[index]) {
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
	resultingStr := string(bytes[index : index+strLen])
	return resultingStr, index + strLen, nil
}

func decodeRType(bytes []byte, index int) (*globals.Rtype, int, error) {
	index++
	var res globals.Rtype
	var key string
	var err error
	readingKey := true
	for index < len(bytes) && bytes[index] != 'e' {
		if readingKey && globals.ByteIsDigit(bytes[index]) {
			key, index, err = decodeStringLiteral(bytes, index)
			if err != nil {
				return nil, 0, err
			}
			readingKey = !readingKey
		} else if readingKey {
			return nil, 0, errors.New(fmt.Sprintf("Dictionary key can be nothing but string. Index: %d, Symbol: '%s'", index, string(bytes[index])))
		} else if !readingKey {
			switch key {
			case "id", "nodes", "token":
				var val string
				val, index, err = decodeStringLiteral(bytes, index)
				if err != nil {
					return nil, 0, err
				}
				switch key {
				case "id":
					res.Id = val
				case "nodes":
					res.Nodes = globals.TakePointer(val)
				case "token":
					res.Token = globals.TakePointer(val)
				}
			case "values":
				var val []string
				val, index, err = decodeStringList(bytes, index)
				if err != nil {
					return nil, 0, err
				}
				res.Values = val

			}
			readingKey = !readingKey
		}
	}
	if index >= len(bytes) {
		return nil, 0, errors.New(fmt.Sprintf("No closing 'e' symbol for main dictionary"))
	}
	return &res, index+1, err
}

func decodeAType(bytes []byte, index int) (*globals.Atype, int, error) {
	var res globals.Atype
	var key string
	var err error
	readingKey := true
	for index < len(bytes) && bytes[index] != 'e' {
		if readingKey && globals.ByteIsDigit(bytes[index]) {
			key, index, err = decodeStringLiteral(bytes, index)
			if err != nil {
				return nil, 0, err
			}
			readingKey = !readingKey
		} else if readingKey {
			return nil, 0, errors.New(fmt.Sprintf("Dictionary key can be nothing but string. Index: %d, Symbol: '%s'", index, string(bytes[index])))
		} else if !readingKey {
			switch key {
			case "id", "target", "info_hash", "token":
				var val string
				val, index, err = decodeStringLiteral(bytes, index)
				if err != nil {
					return nil, 0, err
				}
				switch key {
				case "id":
					res.Id = val
				case "target":
					res.Target = globals.TakePointer(val)
				case "info_hash":
					res.InfoHash = globals.TakePointer(val)
				case "token":
					res.Token = globals.TakePointer(val)
				}
			case "implied_port", "port":
				var port int
				port, index, err = decodeNumber(bytes, index)
				if err != nil {
					return nil, 0, err
				}
				switch key {
				case "implied_port":
					res.ImpliedPort = globals.TakePointer(port)
				case "port":
					res.Port = globals.TakePointer(port)
				}
			}
			readingKey = !readingKey
		}
	}
	if index >= len(bytes) {
		return nil, 0, errors.New(fmt.Sprintf("No closing 'e' symbol for main dictionary"))
	}
	return &res, index+1, err
}

func decodeMainMap(bytes []byte) (*globals.PackageType, error) {
	index := 1
	var res globals.PackageType
	var key string
	var err error
	readingKey := true
	for index < len(bytes) && bytes[index] != 'e' {
		if readingKey && globals.ByteIsDigit(bytes[index]) {
			key, index, err = decodeStringLiteral(bytes, index)
			if err != nil {
				return nil, err
			}
			readingKey = !readingKey
		} else if readingKey {
			return nil, errors.New(fmt.Sprintf("Dictionary key can be nothing but string. Index: %d, Symbol: '%s'", index, string(bytes[index])))
		} else if !readingKey {
			switch key {
			case "e":
				return nil, errors.New(string(bytes[index+2:])) // todo better error handling
			case "t", "y", "q", "v":
				var val string
				val, index, err = decodeStringLiteral(bytes, index)
				if err != nil {
					return nil, err
				}
				switch key {
				case "t":
					res.T = val
				case "y":
					res.Y = val
				case "q":
					res.Q = val
				case "v":
					res.V = val
				}
			case "r":
				var temp *globals.Rtype
				temp, index, err = decodeRType(bytes, index)
				if err != nil {
					return nil, err
				}
				res.R = temp
			case "a":
				var temp *globals.Atype
				temp, index, err = decodeAType(bytes, index)
				if err != nil {
					return nil, err
				}
				res.A = temp
			default:
				return nil, errors.New("unknown key: " + key)
			}
			readingKey = !readingKey
		}
	}
	if index >= len(bytes) {
		return nil, errors.New(fmt.Sprintf("No closing 'e' symbol for main dictionary"))
	}
	return &res, err
}

func Decode(encodedBytes []byte) (*globals.PackageType, error) {
	if len(encodedBytes) > 0 {
		return decodeMainMap(encodedBytes)
	} else {
		return nil, nil
	}
}
