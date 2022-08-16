package main

import (
	"encoding/hex"
	"fmt"

	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

type Opcode byte

const (
	PUSH1   Opcode = 0x60
	PUSH2          = 0x61
	PUSH3          = 0x62
	PUSH32         = 0x7F
	Mstore         = 0x52
	Mstore8        = 0x53
	ADD            = 0x01
	MUL            = 0x02
	SDIV           = 0x05
	EXP            = 0x0a
)

// func extractIntFromByteArray(input []byte, start, length uint64) (len uint64) {
// 	slice := input[start:(start + length)]
// 	strLen := big.NewInt(0).SetBytes(slice).Uint64()

// 	return strLen
// }
func extractHexFromByteArray(input []byte, start, length uint64) uint256.Int {
	slice := input[start:(start + length)]
	sliceString := hex.EncodeToString(slice)
	sliceString = "0x" + sliceString
	result, err := uint256.FromHex(sliceString)
	if err != nil {
		panic(err)
	}
	return *result
}

func calcMemSize64WithUint(off *uint256.Int, length64 uint64) (uint64, bool) {
	// if length is zero, memsize is always zero, regardless of offset
	if length64 == 0 {
		return 0, false
	}
	// Check that offset doesn't overflow
	offset64, overflow := off.Uint64WithOverflow()
	if overflow {
		return 0, true
	}
	val := offset64 + length64
	// if value < either of it's parts, then it overflowed
	return val, val < offset64
}

func toWordSize(size uint64) uint64 {
	if size > math.MaxUint64-31 {
		return math.MaxUint64/32 + 1
	}

	return (size + 31) / 32
}

func solve(input []byte) (uint64, Stack, Memory) {
	length := len(input)
	s := newstack()
	mem := NewMemory()
	integer := new(uint256.Int)
	var gascost uint64 = 0
	var lastgascost uint64 = 0
	for i := 0; i < length; i++ {
		switch Opcode(input[i]) {
		case PUSH1:
			fmt.Println("PUSH1")
			integer.SetBytes(common.RightPadBytes(input[i+1:i+2], 1))
			s.push(integer)
			gascost += 3
			i++
		case PUSH2:
			fmt.Println("PUSH2")
			integer.SetBytes(common.RightPadBytes(input[i+1:i+3], 2))
			s.push(integer)
			gascost += 3
			i += 2
		case PUSH3:
			fmt.Println("PUSH3")
			integer.SetBytes(common.RightPadBytes(input[i+1:i+4], 3))
			s.push(integer)
			gascost += 3
			i += 3
		case PUSH32:
			fmt.Println("PUSH32")
			len := len(input)
			var endIndex int = 0
			if i+33 < len {
				endIndex = i + 33
			} else {
				endIndex = len - 1
			}
			integer.SetBytes(common.RightPadBytes(input[i+1:endIndex], 32))
			s.push(integer)
			gascost += 3
			i += 32
		case Mstore:
			fmt.Println("Mstore")
			offset := s.pop()
			val := s.pop()
			memSize, overflow := calcMemSize64WithUint(&offset, 32)
			if overflow {
				panic("overflow")
			}
			memorySize, overflow := math.SafeMul(toWordSize(memSize), 32)
			if overflow {
				panic("overflow")
			}
			mem.Resize(memorySize)
			mem.Set32(offset.Uint64(), &val)
			w := mem.Len() / 32
			newCost := 3*w + (w * w / 512)
			gascost += (uint64(newCost) - lastgascost) + 3
			lastgascost = uint64(newCost)

		case Mstore8:
			fmt.Println("Mstore8")
			offset := s.pop()
			val := s.pop()
			memSize, overflow := calcMemSize64WithUint(&offset, 1)
			if overflow {
				panic("overflow")
			}
			memorySize, overflow := math.SafeMul(toWordSize(memSize), 32)
			if overflow {
				panic("overflow")
			}
			mem.Resize(memorySize)
			ret := mem.Set(offset.Uint64(), 1, val.Bytes())
			for ret == -1 {
				mem.Resize(uint64(mem.Len()) + 32)
				ret = mem.Set(offset.Uint64(), 1, val.Bytes())
			}
			w := mem.Len() / 32
			newCost := 3*w + (w * w / 512)

			gascost += (uint64(newCost) - lastgascost) + 3
			lastgascost = uint64(newCost)
		case ADD:
			fmt.Println("ADD")
			a := s.pop()
			b := s.pop()
			result := uint256.NewInt(1)
			result.Add(&a, &b)
			s.push(result)
			gascost += 3
		case MUL:
			fmt.Println("MUL")
			a := s.pop()
			b := s.pop()
			result := uint256.NewInt(1)
			result.Mul(&a, &b)
			s.push(result)
			gascost += 5
		case SDIV:
			fmt.Println("SDIV")
			num := s.pop()
			denom := s.pop()
			result := uint256.NewInt(1)
			result.SDiv(&num, &denom)
			s.push(result)
			gascost += 5
		case EXP:
			fmt.Println("EXP")
			base := s.pop()
			exponent := s.pop()
			result := uint256.NewInt(1)
			expLength := (exponent.BitLen() + 7) / 8
			result.Exp(&base, &exponent)
			s.push(result)
			gascost += uint64(expLength*50) + 10
		}
	}
	return gascost, *s, *mem
}
func findHash(memory []byte) common.Hash {
	hash := crypto.Keccak256Hash(memory)
	return hash
}
func findSolution(inputString string) (uint64, Stack, Memory, string) {

	byteArray, err := hex.DecodeString(inputString)
	if err != nil {
		fmt.Println("Unable to convert hex to byte. ", err)
	}
	result, s, m := solve(byteArray)
	hash := findHash(m.store)
	return result, s, m, hash.String()
}

func main() {
	input := os.Args[1:]
	inputString := strings.Join(input, "")
	gas, stack, mem, hash := findSolution(inputString)
	fmt.Println("gas  : ", gas)
	fmt.Println("hash : ", hash)
	fmt.Println("stack : ", stack)
	fmt.Println("memory : ", mem)

}
