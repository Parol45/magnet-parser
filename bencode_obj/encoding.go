package bencode_obj

import (
	"errors"
	"magnet-parser/globals"
	"strconv"
)

func encodeStringLiteral(literal string) []byte {
	var resultingBytes, temp []byte
	temp = []byte(literal)
	resultingBytes = append(resultingBytes, strconv.Itoa(len(temp))...)
	resultingBytes = append(resultingBytes, ':')
	resultingBytes = append(resultingBytes, temp...)
	return resultingBytes
}

func Encode(obj *globals.PackageType) (encodedBytes []byte, err error) {
	if obj == nil {
		return
	} else if obj.A == nil && obj.R == nil {
		err = errors.New("package doesn't contain nor arguments nor response fields")
		return
	}

	encodedBytes = append(encodedBytes, 'd')
	if obj.A != nil {
		encodedBytes = append(encodedBytes, []byte("1:ad2:id")...)
		encodedBytes = append(encodedBytes, encodeStringLiteral(obj.A.Id)...)
		// todo....
		encodedBytes = append(encodedBytes, 'e')
	} else if obj.R != nil {

	}
	encodedBytes = append(encodedBytes, []byte("1:q")...)
	encodedBytes = append(encodedBytes, encodeStringLiteral(obj.Q)...)

	encodedBytes = append(encodedBytes, []byte("1:t")...)
	encodedBytes = append(encodedBytes, encodeStringLiteral(obj.T)...)

	encodedBytes = append(encodedBytes, []byte("1:y")...)
	encodedBytes = append(encodedBytes, encodeStringLiteral(obj.Y)...)


	encodedBytes = append(encodedBytes, 'e')
	str := string(encodedBytes)
	println(str[12:32])
	return
}
