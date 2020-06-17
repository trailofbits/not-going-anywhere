package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	_, err := os.Create("./proof")
	if err != nil {
		buf, err := ioutil.ReadFile("./proof")
		if err != nil {
			fmt.Println("failed to read file")
		}
		fmt.Printf("%s", buf)
		// fmt.Printf("the file got made at %d and is %s", time.Now() fs.Name())
	}
	fmt.Println(err)
}
