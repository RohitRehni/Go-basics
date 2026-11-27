package main

import (
	"fmt"
	"strings"
)

func main() {

	fmt.Println(romanToInt("III"))
	fmt.Println(romanToInt("LVIII"))
	fmt.Println(romanToInt("MCMXCIV"))
	fmt.Println(romanToInt("MCMLXXXV"))
	fmt.Println(romanToInt("MCMXC"))
	fmt.Println(romanToInt("MMVI"))
}

var romanNumerals = map[string]int{
	"I": 1,
	"V": 5,
	"X": 10,
	"L": 50,
	"C": 100,
	"D": 500,
	"M": 1000,
}

var romanOrder = map[string]int{
	"I": 1,
	"V": 2,
	"X": 3,
	"L": 4,
	"C": 5,
	"D": 6,
	"M": 7,
}

func romanToInt(s string) int {

	var int_array = strings.Split(s, "")

	var sum int
	for i := 0; i < len(int_array); i++ {
		var Number int
		var current_ele string = int_array[i]

		var prev_ele string
		if i > 0 {
			prev_ele = int_array[i-1]
		}

		if romanOrder[prev_ele] < romanOrder[current_ele] {
			Number = romanNumerals[current_ele] - romanNumerals[prev_ele]
			sum = sum - romanNumerals[prev_ele]
		} else {
			Number = romanNumerals[current_ele]
		}

		sum = sum + Number
	}

	return sum
}
