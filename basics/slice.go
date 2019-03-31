package main

import "fmt"

func p(s []int) {
	fmt.Printf("%#v, len=%d, cap=%d\n", s, len(s), cap(s))
}

func main() {
	a := make([]int, 5, 10)
	p(a) // => []int{0, 0, 0, 0, 0}, len=5, cap=10

	b := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	p(b) // => []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, len=10, cap=10

	c := b[2:5]
	p(c) // => []int{2, 3, 4}, len=3, cap=8

	d := b[2:5:8]
	p(d) // => []int{2, 3, 4}, len=3, cap=6

	c = append(c, 10, 11, 12)
	p(b) // => []int{0, 1, 2, 3, 4, 10, 11, 12, 8, 9}, len=10, cap=10
	p(c) // => []int{2, 3, 4, 10, 11, 12}, len=6, cap=8

	e := []int{0, 1, 2, 3, 4, 5}
	f := []int{6, 7, 8}
	copy(f, e)
	p(f) // => []int{0, 1, 2}, len=3, cap=3
}
