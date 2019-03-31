package main

import "fmt"

func variadicFunction(s ...int) {
	for _, v := range s {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}

func main() {
	variadicFunction([]int{1, 2, 3, 4, 5}...)
	// => 1 2 3 4 5
}
