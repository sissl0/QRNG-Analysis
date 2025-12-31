package src

import (
	"fmt"
	"os"
)

func getHistogram(filename string) {

	// Read entire file
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Count array for all 256 possible 8-bit patterns
	var counts [256]uint64

	var window uint8 = 0
	var bitsSeen int = 0

	for _, b := range data {
		// Process bits MSB -> LSB
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1

			// Shift left and insert new bit
			window = (window << 1) | uint8(bit)

			if bitsSeen >= 7 {
				counts[window]++
			}

			bitsSeen++
		}
	}

	// Print results
	for i := 0; i < 256; i++ {
		fmt.Printf("%08b : %d\n", i, counts[i])
	}
}
