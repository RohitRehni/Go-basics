package main

import (
	"fmt"
)

func main() {

	var str = []string{"flower", "flow", "flight"}
	fmt.Println("1 :::::::::::::::::::::::::::: " + longestCommonPrefix(str))

	str = []string{"dog", "racecar", "car"}
	fmt.Println("2 :::::::::::::::::::::::::::: " + longestCommonPrefix(str))

	str = []string{"flower", "flower", "flower"}
	fmt.Println("3 :::::::::::::::::::::::::::: " + longestCommonPrefix(str))

	str = []string{"flower", "fkow"}
	fmt.Println("4 :::::::::::::::::::::::::::: " + longestCommonPrefix(str))

}

func longestCommonPrefix(strs []string) string {

	if len(strs) == 0 {
		return ""
	}

	prefix := strs[0]

	for _, s := range strs[1:] {

		for len(prefix) > 0 {

			if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
				break
			}

			prefix = prefix[:len(prefix)-1]
		}

		if prefix == "" {
			return ""
		}
	}

	return prefix
}

func longest_approch_common_prefix(strs []string) string {

	// find, string which has the least lenght
	// iterate, to that string to find the all prefixes that's possible
	// now, find these iterations of prefixes in other strings and create a string map count
	// get, the highest string map count and that is the prefix

	fmt.Println(strs)

	var str_length = len(strs)
	var min_string string = strs[0]
	var min_length int = len(min_string)

	// for _, value := range strs {
	// 	var length int = len(value)
	// 	if (length < min_length ) {
	// 		min_length = length
	// 		min_string = value
	// 	}
	// }

	var prefix_str_cnt = make([]string, min_length)

	for i := 1; i <= min_length; i++ {
		var slice string = min_string[0:i]
		prefix_str_cnt[i-1] = slice
	}

	prefix_count_map := make(map[string]int)

	for _, value := range strs {
		for _, prefix_str := range prefix_str_cnt {

			var prefix_len int = len(prefix_str)

			if prefix_len > len(value) {
				prefix_len = len(value)
			}

			if prefix_len < 1 {
				continue
			}

			var prefix_str_args string = value[0:prefix_len]
			if prefix_str_args == prefix_str {
				count, exists := prefix_count_map[prefix_str]
				if exists {
					prefix_count_map[prefix_str] = count + 1
				} else {
					prefix_count_map[prefix_str] = 1
				}
			} else {
				break
			}
		}
	}

	fmt.Println("Prefix String Maps :: ")
	fmt.Println(prefix_count_map)

	var max_common_string_length int = 0
	var max_common_string string

	fmt.Println("String Length :: ")
	fmt.Println(str_length)

	var found int = 0
	for key, value := range prefix_count_map {

		if found == 1 {

			if value == str_length {
				if max_common_string_length < len(key) {
					max_common_string = key
					max_common_string_length = len(key)
				}
			}

		} else {

			if value == str_length {
				found = 1
				max_common_string = key
				max_common_string_length = len(key)
			}
		}
	}

	return max_common_string
}
