package main

import "fmt"

func main() {
	a := []byte("ba")

	a1 := append(a, 'd')
	a2 := append(a, 'g')

	fmt.Println(string(a1)) // bag
	fmt.Println(string(a2)) // bag
}
