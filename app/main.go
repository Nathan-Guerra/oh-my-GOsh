package main

import (
	"fmt"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	var command string
	// TODO: Uncomment the code below to pass the first stage
	for {
		fmt.Print("$ ")
		fmt.Scanf("%s", &command)
		fmt.Printf("%s: command not found\n", command)
	}
}
