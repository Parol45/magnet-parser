package bencode_obj

import (
	"errors"
	"magnet-parser/globals"
	"strconv"
	"strings"
)

func encodeStringLiteral(literal string) []byte {
	var resultingBytes, temp []byte
	temp = []byte(literal)
	resultingBytes = append(resultingBytes, strconv.Itoa(len(temp))...)
	resultingBytes = append(resultingBytes, ':')
	resultingBytes = append(resultingBytes, temp...)
	return resultingBytes
}

func encodeNumber(number int) []byte {
	builder := strings.Builder{}
	builder.WriteByte('i')
	builder.WriteString(strconv.Itoa(number))
	builder.WriteByte('e')
	return []byte(builder.String())
}

func encodeStringList(list []string) []byte {
	var resultingBytes []byte
	resultingBytes = append(resultingBytes, 'l')
	for _, literal := range list {
		resultingBytes = append(resultingBytes, encodeStringLiteral(literal)...)
	}
	resultingBytes = append(resultingBytes, 'e')
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
	// ------------------------------------------------------------------------------------------------------
	if obj.A != nil {
		encodedBytes = append(encodedBytes, []byte("1:ad2:id")...)
		encodedBytes = append(encodedBytes, encodeStringLiteral(obj.A.Id)...)

		if obj.A.ImpliedPort != nil {
			encodedBytes = append(encodedBytes, []byte("12:implied_port")...)
			encodedBytes = append(encodedBytes, encodeNumber(*obj.A.ImpliedPort)...)
		}

		if obj.A.InfoHash != nil {
			encodedBytes = append(encodedBytes, []byte("9:info_hash")...)
			encodedBytes = append(encodedBytes, encodeStringLiteral(*obj.A.InfoHash)...)
		}

		if obj.A.Port != nil {
			encodedBytes = append(encodedBytes, []byte("4:port")...)
			encodedBytes = append(encodedBytes, encodeNumber(*obj.A.Port)...)
		}

		if obj.A.Target != nil {
			encodedBytes = append(encodedBytes, []byte("6:target")...)
			encodedBytes = append(encodedBytes, encodeStringLiteral(*obj.A.Target)...)
		}

		if obj.A.Token != nil {
			encodedBytes = append(encodedBytes, []byte("5:token")...)
			encodedBytes = append(encodedBytes, encodeStringLiteral(*obj.A.Token)...)
		}
		encodedBytes = append(encodedBytes, 'e')
	}
	// ------------------------------------------------------------------------------------------------------
	encodedBytes = append(encodedBytes, []byte("1:q")...)
	encodedBytes = append(encodedBytes, encodeStringLiteral(obj.Q)...)
	// ------------------------------------------------------------------------------------------------------
	if obj.R != nil {
		encodedBytes = append(encodedBytes, []byte("1:rd2:id")...)
		encodedBytes = append(encodedBytes, encodeStringLiteral(obj.R.Id)...)
		if obj.R.Nodes != nil {
			encodedBytes = append(encodedBytes, []byte("5:nodes")...)
			encodedBytes = append(encodedBytes, encodeStringLiteral(*obj.R.Nodes)...)
		}

		if obj.R.Token != nil {
			encodedBytes = append(encodedBytes, []byte("5:token")...)
			encodedBytes = append(encodedBytes, encodeStringLiteral(*obj.R.Token)...)
		}

		if obj.R.Values != nil {
			encodedBytes = append(encodedBytes, []byte("6:values")...)
			encodedBytes = append(encodedBytes, encodeStringList(obj.R.Values)...)
		}
		encodedBytes = append(encodedBytes, 'e')
	}
	// ------------------------------------------------------------------------------------------------------
	encodedBytes = append(encodedBytes, []byte("1:t")...)
	encodedBytes = append(encodedBytes, encodeStringLiteral(obj.T)...)
	// ------------------------------------------------------------------------------------------------------
	encodedBytes = append(encodedBytes, []byte("1:y")...)
	encodedBytes = append(encodedBytes, encodeStringLiteral(obj.Y)...)
	// ------------------------------------------------------------------------------------------------------
	encodedBytes = append(encodedBytes, 'e')
	return
}
