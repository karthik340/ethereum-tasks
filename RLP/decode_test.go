package main

import (
	"testing"
)

func TestDecode(t *testing.T) {

	for _, tc := range []struct {
		test             string
		expectedSolution string
		isErrorExpected  bool
	}{
		{
			"ed90416e616e746861204b726973686e616e8d526168756c204c656e6b616c618d47616e65736820507261736164",
			"List { \nAnantha Krishnan,Rahul Lenkala,Ganesh Prasad \n}",
			false,
		},
		{
			"e5922034342e38313538393735343033373334319132302e3435343733343334343535353435",
			"List { \n 44.81589754037341,20.45473434455545 \n}",
			false,
		},
		{
			"3f",
			"?",
			false,
		},
		{
			"b8686162636465666768696a6b6c6d6e6f707172747374757678797a6162636465666768696a6b6c6d6e6f707172747374757678797a6162636465666768696a6b6c6d6e6f707172747374757678797a6162636465666768696a6b6c6d6e6f707172747374757678797a",
			"abcdefghijklmnopqrtstuvxyzabcdefghijklmnopqrtstuvxyzabcdefghijklmnopqrtstuvxyzabcdefghijklmnopqrtstuvxyz",
			false,
		},
		{
			"f8659862636465666768696a6b6c6d6e6f7071727374757678797a9862636465666768696a6b6c6d6e6f7071727374757678797a9862636465666768696a6b6c6d6e6f7071727374757678797a996162636465666768696a6b6c6d6e6f7071727374757678797a",
			"List { \nbcdefghijklmnopqrstuvxyz,bcdefghijklmnopqrstuvxyz,bcdefghijklmnopqrstuvxyz,abcdefghijklmnopqrstuvxyz \n}",
			false,
		},
		{
			"ed",
			"",
			true,
		},

		{
			"c5c4c3c2c180",
			"List { \nList { \nList { \nList { \nList { \n \n} \n} \n} \n} \n}",
			false,
		},
	} {
		result := findSolution(tc.test)
		if tc.expectedSolution != result {
			if !tc.isErrorExpected {
				t.Errorf("incorrect expected : %s actual %s ", tc.expectedSolution, result)
			} else {
				t.Errorf("incorrect expected : %s actual %s ", tc.expectedSolution, result)
			}
		}
	}

}
