package main

import "fmt"

func map_print(m map[string]int) {
	fmt.Printf("\n")

	for k, v := range m{
		fmt.Printf("M: %s D: %d\n", k, v)
	}
}

func main() {
	monthdays := map[string]int{
		"Jan": 31, "Feb": 28, "Mar": 31,
		"Apr": 30, "May": 31, "Jun": 30,
		"Jul": 31, "Aug": 31, "Sep": 30,
		"Oct": 31, "Nov": 30, "Dec": 31,
	}

	map_print(monthdays)

	delete(monthdays, "Jan")

	map_print(monthdays)

	monthdays["Jan"] = 31

	map_print(monthdays)

	delete(monthdays, "Jan")

	map_print(monthdays)
}

