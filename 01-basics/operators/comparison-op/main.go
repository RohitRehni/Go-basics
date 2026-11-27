package main

import (
	"fmt"
	"math"
	"reflect"
	"slices"
)

func main() {
	fmt.Println("Operations for Comparision Op")

	//  Compile-time error.
	// Go does not treat true as 1 or false as 0.
	// true == 1 (mismatched types untyped bool and untyped int)
	// fmt.Println(true == 1)

	var str1 string = "String1"
	var str2 string = "String2"

	if str1 == str2 {
		fmt.Print("String Matching : TRUE" + "\n")
	} else {
		fmt.Print("String Matching : FALSE" + "\n")
	}

	var num1 int = 10
	var num2 int = 200

	// invalid operation: num1 == str2 (mismatched types int and string)
	// // if (num1 == str2) {
	// // 	fmt.Print("TRUE" + "\n")
	// // }

	if num1 != num2 {
		fmt.Print("Number Matching : TRUE" + "\n")
	}

	// fmt.Println(5 == 5.0) // Compile-time error.
	// Go does not compare an integer and a floating-point value directly, even when their values are equal.
	fmt.Println(float64(5) == 5.0)

	// Floating-point values may not be represented exactly in binary.
	fmt.Println("float point comapre : ", 0.1+0.2 == 0.3) // false
	fmt.Println("float point comapre using math : ", math.Abs((0.1+0.2)-0.3) < 1e-9)

	x_nan := math.NaN()

	// NaN is not equal to itself.
	fmt.Println(x_nan == x_nan)
	fmt.Println(x_nan != x_nan)

	// Strings are compared lexicographically, byte by byte. "10" comes before "2" because '1' < '2'.
	fmt.Println(" Strings are compared lexicographically : ", "apple" < "banana")
	fmt.Println(" Integers as Strings are compared lexicographically : ", "10" < "2")

	// String comparisons are case-sensitive.
	fmt.Println("String comparisons are case-sensitive.", "Go" == "go")

	// An interface contains both a dynamic type and a value:
	var p *int = nil
	var x any = p

	// The interface itself is not nil.
	fmt.Println(x == nil)

	// Dynamic type have same type and Value.
	var a any = 10
	var b any = 10
	fmt.Println(a == b)

	var a_slice any = []int{1, 2}
	var b_slice any = []int{1, 2}

	// panic: runtime error: comparing uncomparable type []int
	// fmt.Println("dynamic slice : " , a_slice == b_slice)

	fmt.Println("dynamic slice comapre using reflect: ", reflect.DeepEqual(a_slice, b_slice))
	fmt.Println("dynamic slice comapre using slices Equal : ", slices.Equal(a_slice.([]int), b_slice.([]int)))

	// check if slice is empty or nil
	a_slice_sh := []int{1, 2}
	fmt.Println(a_slice_sh == nil)

	// a is a nil slice. b is an initialized, empty slice.
	// Both have length zero, but they are not both nil.
	var a_declared []int
	b_intialized := []int{}
	fmt.Println("Declare Slice compare to nil", a_declared == nil)
	fmt.Println("Initialized Slice compare to nil", b_intialized == nil)

	// Function values cannot be compared to each other.
	// Maps cannot be compared with each other. They can only be compared with nil:

	// Structs are comparable if all their fields are comparable.
	type Point struct {
		X int
		Y int
	}

	a_struct := Point{1, 2}
	b_struct := Point{1, 2}
	fmt.Println("Comapring 2 Struct contaning ints : ", a_struct == b_struct)

	//  The struct contains a slice, so the struct itself is not comparable.
	type Data struct {
		Values []int
	}

	// invalid operation: a_slice_struct == b_slice_struct (struct containing []int cannot be compared)
	// a_slice_struct := Data{Values: []int{1}}
	// b_slice_struct := Data{Values: []int{1}}

	// fmt.Println("Comapring 2 Struct contaning slices :", a_slice_struct == b_slice_struct)

	type UserID int

	// var a_int int = 10
	// var b_type_user UserID = 10

	// cannot convert b (variable of interface type any) to type int: need type assertion
	// fmt.Println(a_int == b_type_user)
	// int and UserID are different named types, even though UserID has an underlying type of int.
	// fmt.Println("Compare int to a User Defined Type : ", a_int == int(b_type_user))

	// The pointers point to different variables, even though the variables contain equal values.
	x_int := 10
	y_int := 10
	fmt.Println("compare Pointers : ", &x_int == &y_int)

	x_int_value := 10
	int_pointer := &x_int_value

	fmt.Println("compare Pointer and Int type : ", *int_pointer == x_int_value)

}
