package main

import (
	"fmt"
	"strconv"
)

func main() {

	var y bool = is_palindrome(101)
	if y {
		fmt.Println("is_palindrome is true")
	} else {
		fmt.Println("is_palindrome is false")
	}
}

func is_palindrome(x int) bool {

	var reverse int = 0
	var temp int = x
	var neg_sign string = ""
	if x < 0 {
		neg_sign = "-"
		temp = temp * -1
	}

	var count int = 0

	for temp != 0 {

		var a int = (reverse * 10)
		var b int = (temp % 10)

		reverse = a + b
		temp = temp / 10

		count++
	}

	// For keeping the Reverse Value have zero's if missed any
	reverse_str := fmt.Sprintf(fmt.Sprintf("%%0%dd%v", count, neg_sign), reverse)
	var x_str = strconv.Itoa(x)
	if x_str == reverse_str {
		return true
	}

	return false
}
