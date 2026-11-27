package main

import (
	"fmt"
)

func main() {

	fmt.Println("Operations via arrays")

	// Initialization
	var grades [3]int = [3]int{10, 20, 30}
	fmt.Println("grades ::", grades)

	grades_shorthand := [3]int{10, 20, 30}
	fmt.Println("grades_shorthand ::", grades_shorthand)

	grade_elipses := [...]int{10, 20, 30, 90, 80}
	fmt.Println("grade_elipses ::", grade_elipses)

	// Length
	fmt.Println("Length of grade_elipses :", len(grade_elipses))

	// Index of last element : Length - 1
	// First Index : 0
	fmt.Println("grades Index 0 ::", grades[0])
	fmt.Println("grades Index Last Index ::", len(grades)-1, ", Values as ::", grades[len(grades)-1])

	// Out of Bound : invalid array index, if the index is past element array

	grades_shorthand[0] = 30
	fmt.Println("grades_shorthand Index 0 after updating::", grades_shorthand[0])

	// Looping
	for i := 0; i < len(grade_elipses); i++ {
		fmt.Println("using for Loop call Index : ", i, ", value :", grade_elipses[i])
	}

	// Range Keyword
	for index, element := range grade_elipses {
		fmt.Println("using Range keyword call Index : ", index, "=> value :", element)
	}

	// Array of array or 2d array or 3d array
	var array_d [2][2]int = [2][2]int{{1, 2}, {2, 2}}
	for index, element := range array_d {
		fmt.Println("2D Array Index : ", index, ", element :", element)

		for inner_index, inner_element := range element {
			fmt.Println("2D Array inner_index : ", inner_index, ", element :", inner_element)
		}
	}

}
