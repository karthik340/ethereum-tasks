package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
)

const CHAR = 1
const STRING = 2
const LIST = 3

func extractIntFromByteArray(input []byte, start, length uint64) (len uint64) {
	slice := input[start:(start + length)]
	strLen := big.NewInt(0).SetBytes(slice).Uint64()

	return strLen
}
func extractStringFromByteArray(input []byte, start, length uint64) (res string) {
	result := string(input)[start:(start + length)]
	return result
}
func identifyType(input []byte) (offset uint64, length uint64, Type uint64, err error) {
	inputLength := len(input)
	if inputLength == 0 {
		return 0, 0, 0, errors.New("input is null")
	}
	firstByte := input[0]

	if firstByte >= byte(0x00) && firstByte <= byte(0x7f) {
		return 0, 1, CHAR, nil
	} else if firstByte >= byte(0x80) && firstByte <= byte(0xb7) && inputLength > int(firstByte-byte(0x80)) {
		strLen := firstByte - byte(0x80)
		return 1, uint64(strLen), STRING, nil
	} else if firstByte >= byte(0xb8) && firstByte <= byte(0xbf) && inputLength > int(firstByte-byte(0xb7)) && inputLength > int(firstByte-byte(0xb7))+int(extractIntFromByteArray(input, 1, uint64(firstByte-byte(0xb7)))) {
		secondPartLength := uint64(firstByte - byte(0xb7))
		strLen := extractIntFromByteArray(input, 1, secondPartLength)
		return 1 + secondPartLength, strLen, STRING, nil
	} else if firstByte >= byte(0xc0) && firstByte <= byte(0xf7) && inputLength > int(firstByte-byte(0xc0)) {
		strLen := firstByte - byte(0xc0)
		return 1, uint64(strLen), LIST, nil
	} else if firstByte >= byte(0xf8) && firstByte <= byte(0xff) && inputLength > int(firstByte-byte(0xf7))+int(extractIntFromByteArray(input, 1, uint64(firstByte-byte(0xf7)))) {
		secondPartLength := uint64(firstByte - byte(0xf7))
		strLen := extractIntFromByteArray(input, 1, secondPartLength)
		return 1 + secondPartLength, strLen, LIST, nil
	}
	return 0, 0, 0, errors.New("unkown error")
}

func extractStringsAndList(input []byte) string {
	if len(input) == 0 {
		return ""
	}

	var result, result1 string

	offset, length, Type, err := identifyType(input)

	if err == nil {
		if Type == CHAR {
			result = extractStringFromByteArray(input, offset, length)
			result1 = extractStringsAndList(input[offset+length:])
		} else if Type == STRING {
			result = extractStringFromByteArray(input, offset, length)
			result1 = extractStringsAndList(input[offset+length:])
		} else if Type == LIST {
			result1 = "List { \n" + extractStringsAndList(input[offset:]) + " \n}"
		}
		if result == "" {
			return result1
		}
		if result1 == "" {
			return result
		}
		result = result + "," + result1
	} else {
		fmt.Println("unable to decode")
	}

	return result
}

func findSolution(inputString string) string {
	byteArray, err := hex.DecodeString(inputString)
	if err != nil {
		fmt.Println("Unable to convert hex to byte. ", err)
	}
	result := extractStringsAndList(byteArray)
	return result
}
func main() {
	input := os.Args[1:]
	inputString := strings.Join(input, "")

	fmt.Println(findSolution(inputString))
}
