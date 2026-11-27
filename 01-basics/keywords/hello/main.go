package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello Go")

	var name = "Starting of Learning Basics"
	fmt.Println("\n" + name)

	// Using fmt.Sprintf which return a concatenated string
	var s1 = "Repeat after "
	var s2 = "me"
	fmt.Println(fmt.Sprintf("%s%s", s1, s2))

	// reverse a sign to back for a gieven number and Zero for
	reverse_marker := fmt.Sprintf(fmt.Sprintf("%%0%dd%v", 4, "+"), 34)
	fmt.Println(reverse_marker)

}
