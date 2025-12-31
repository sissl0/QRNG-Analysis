package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: go run qrng.go <mode>")
		return
	}
	mode := args[1]
	switch mode {
	case "histogram":
		if len(args) != 3 {
			fmt.Println("Usage: go run qrng.go histogram <filename>")
			return
		}
		getHistogram(args[2])
	default:
		fmt.Println("Unknown mode:", mode)
	}
}
