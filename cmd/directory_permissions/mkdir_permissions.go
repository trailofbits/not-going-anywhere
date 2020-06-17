// This example does not work on windows machines
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// hypothetically an attacker does this without your knowledge
	err := os.MkdirAll("./proof", 0777)

	// you try to make a directory with appropriate permissions
	err = os.MkdirAll("./proof", 0644)
	fmt.Println("No error occurs even though the directory with overly permissive settings still exists: ", err)
	if err = ll(); err != nil {
		fmt.Printf("Error when calling ls -la: %s\n", err)
	}

	info, _ := os.Stat("./proof")
	// Verify if what you made matches the permissions you intended before using with info.Mode()
	if info.Mode() != 0644 {
		fmt.Println("fixing permissions after checking...")
		if err = os.Chmod("proof", 0644); err != nil {
			fmt.Printf("Failed to chmod permissions on proof: %s\n", err)
		}
	}

	if err = ll(); err != nil {
		fmt.Printf("Error when calling ls -la: %s\n", err)
	}
}

func ll() error {
	fmt.Println("$ ls -la")
	ls := exec.Command("ls", "-la")
	ls.Stdout = os.Stdout
	return ls.Run()
}
