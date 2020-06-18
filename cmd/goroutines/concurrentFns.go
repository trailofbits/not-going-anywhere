package main

import (
	"fmt"
	"sync"
)

func main() {
	ConcurrentFunctions(func1, func2)
}

func ConcurrentFunctions(fns ...func()) {
	var wg sync.WaitGroup
	for _, fn := range fns {
		wg.Add(1)
		go func() {
			fn()
			wg.Done()
		}()
	}

	wg.Wait()
}

func func1() {
	fmt.Println("I am function func1")
}

func func2() {
	fmt.Println("I am function func2")
}
