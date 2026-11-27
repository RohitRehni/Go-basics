package main

import (
	"fmt"
)

// := can't be used outside as the global operator or outside a function

func main() {

	fmt.Println("Operations for Arthmatic Op")

	// General purpose
	// * Integer versus floating-point division
	// * Remainder with negative numbers
	// * Explicit type conversions
	// * Constant evaluation
	// * Integer overflow
	// * Unsigned wraparound
	// * Operator precedence
	// * ++ and -- restrictions
	// * Division by zero
	// * Untyped constants and arbitrary-precision values

	// 4 * 6 = 24 ; Remainder is 3
	fmt.Println(27 % 4)

	i := 10
	i++
	fmt.Println(i)

	// syntax error: unexpected ++ at end of statement
	// :: Go also does not support prefix increment So,
	// y := i++
	// fmt.Println(y)

	// result is floating point value
	fmt.Println(5.0 / 2)

	// - Remainder with negative numbers

	// The remainder has the same sign as the left-hand operand.
	fmt.Println(-7 / 3) // -2
	fmt.Println(-7 % 3) // -1

	// - BODMAS

	// First multiply then add
	// Multiplication has higher precedence than addition.
	// BODMAS (Brackets Orders Division Multiplication Addition Subtraction)
	fmt.Println(2 + 3*4)

	// Operator with same precedence operate from Left to Right , So : 5
	fmt.Println(10 - 3 - 2)

	// left to right 2 / 3 is integer valuation so becomes 0
	fmt.Println(2 / 3 * 3)

	// - Explict Conversion

	var x_unsigned uint = 5
	var y_signed int = 3

	// Go does not allow arithmetic between uint and int without conversion.
	// You must explicitly convert:
	fmt.Println(int(x_unsigned) - y_signed)

	// converting a negative int to uint can produce a very large unsigned value.
	fmt.Println(x_unsigned - uint(y_signed))

	// - Value Overflow typed Memory scope

	var unsigned_int uint8 = 255
	// Unsigned integer arithmetic wraps around when it overflows.
	unsigned_int++
	fmt.Println(unsigned_int)
	unsigned_int--
	fmt.Println(unsigned_int)

	// cannot use 128 (untyped int constant) as int8 value in constant declaration (overflows)
	// const x int8 = 128   // Use 127 as this will stop overflow and const can store one value without overflow
	// const x uint8 = 256  // Use 255

	// - Division by zero

	//  0.0 / 0.0 give error as invalid operation: division by zero but, float gives NaN
	x_float := 0.00
	y_float := 0.00
	fmt.Println(x_float / y_float) // NaN (not a number)

	// This causes a compile-time error because x is an integer and 0.0 is a floating-point value.
	// Go does not implicitly convert x to float64.
	x_int := 10
	// fmt.Println(x_int / 0.0)
	fmt.Println(float64(x_int) / 0.0)

	// - Value out of typed memory scope

	var x int8 = -128
	// No change as 128 can't be represented in int8
	// int8 = -128 to 127
	fmt.Println(-x)

	// this 4611686018427387904 value can fit but (1 << 100) can't
	const x_const int = 1 << 62
	fmt.Println(x_const)

	var x_bitwise, y_bitwise int = 100, 90
	fmt.Println(x_bitwise & y_bitwise)
	fmt.Println(x_bitwise | y_bitwise)

	fmt.Println((x_bitwise + y_bitwise) >> 2)

}
