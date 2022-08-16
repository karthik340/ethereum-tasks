package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDec(t *testing.T) {

	for _, tc := range []struct {
		name             string
		test             string
		expectedSolution string
		expectedGas      uint64
		isErrorExpected  bool
	}{
		{
			"PUSH1",
			"60FF6000",
			"[[255 0 0 0] [0 0 0 0]]",
			6,
			false,
		},
		{
			"PUSH2",
			"61FFFF610000",
			"[[255 0 0 0] [0 0 0 0]]",
			6,
			false,
		},
		{
			"PUSH3",
			"62FFFFFF62000000",
			"[[255 0 0 0] [0 0 0 0]]",
			6,
			false,
		},
		{
			"PUSH32",
			"7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF7F0000000000000000000000000000000000000000000000000000000000000000",
			"[[255 0 0 0] [0 0 0 0]]",
			6,
			false,
		},
		{
			"ADD",
			"7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF7F000000000000000000000000000000000000000000000000000000000000000101",
			"[[255 0 0 0] [0 0 0 0]]",
			9,
			false,
		},
		{
			"ADD PUSH1",
			"600A600A01",
			"[[255 0 0 0] [0 0 0 0]]",
			9,
			false,
		},
		{
			"MUL PUSH1",
			"600A600A02",
			"[[255 0 0 0] [0 0 0 0]]",
			11,
			false,
		},
		{
			"SDIV PUSH1",
			"600A600A05",
			"[[255 0 0 0] [0 0 0 0]]",
			11,
			false,
		},
		{
			"EXP PUSH1",
			"600A600A0A",
			"[[255 0 0 0] [0 0 0 0]]",
			66,
			false,
		},
		{
			"MSTORE PUSH1",
			"60FF60005260FF600152",
			"[[255 0 0 0] [0 0 0 0]]",
			24,
			false,
		},
		{
			"MSTORE8 PUSH2 PUSH1",
			"61FFFF60005360FF600153",
			"[[255 0 0 0] [0 0 0 0]]",
			21,
			false,
		},
		{
			"Problem 1",
			"60016020526002606452600361ff0052600362ffffff526005601053",
			"[[255 0 0 0] [0 0 0 0]]",
			538445872,
			false,
		},
		{
			"Problem 2",
			"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00016000527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00026020527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff05604052",
			"[[255 0 0 0] [0 0 0 0]]",
			58,
			false,
		},
		{
			"Problem 3",
			"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000a6000527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000a6020527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0a604052",
			"[[255 0 0 0] [0 0 0 0]]",
			4875,
			false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gas, s, m, _ := findSolution(tc.test)
			fmt.Println("gas", gas)
			require.Equal(t, tc.expectedGas, gas)

			fmt.Println("stack   - ", s)
			fmt.Println("memory   - ", m)
		})

	}

}

func TestDecodeWithHash(t *testing.T) {

	for _, tc := range []struct {
		name             string
		test             string
		expectedSolution string
		expectedGas      uint64
		isErrorExpected  bool
		expectedHash     string
	}{
		{
			"Problem 1",
			"60016020526002606452600361ff0052600362ffffff526005601053",
			"[[255 0 0 0] [0 0 0 0]]",
			538445872,
			false,
			"0xab2744998886b708acadc0a32428d0aa1953e83924383d21c6de5dac852ccbcc",
		},
		{
			"Problem 2",
			"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00016000527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00026020527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff05604052",
			"[[255 0 0 0] [0 0 0 0]]",
			58,
			false,
			"0xb9a07dba38aa24923a611fced9d2eede3bfbfa281e5e498d60f4bd99e5ce6a15",
		},
		{
			"Problem 3",
			"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000a6000527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000a6020527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0a604052",
			"[[255 0 0 0] [0 0 0 0]]",
			4875,
			false,
			"0xafe1e714d2cd3ed5b0fa0a04ee95cd564b955ab8661c5665588758b48b66e263",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gas, _, _, hash := findSolution(tc.test)
			// fmt.Println("gas", 21000+gas)
			require.Equal(t, tc.expectedGas, gas)

			// fmt.Println("stack   - ", s)
			// fmt.Println("memory   - ", m)
			require.Equal(t, tc.expectedHash, hash)
		})

	}

}
