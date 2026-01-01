package main

import (
	"fmt"
	"math"
	"os"

	"qrnganalysis/src/cmd"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: go run qrng.go <mode>")
		return
	}

	mode := args[1]

	switch mode {

	case "histochi":
		if len(args) != 4 {
			fmt.Println("Usage: go run qrng.go histochi <filename> <outdir>")
			return
		}

		filename := args[2]
		outDir := args[3]

		bits := cmd.LoadBits(filename, cmd.LittleEndian)

		const SketchBits = 16

		for windowBits := 4; windowBits <= 8192; windowBits++ {
			for shift := 0; shift < windowBits; shift++ {

				sketch := cmd.NewSketch(SketchBits)
				cmd.FillSketch(bits, windowBits, shift, sketch)

				chi := cmd.ChiSquareSketch(sketch)
				h := cmd.EntropySketch(sketch)
				hMax := math.Log2(float64(len(sketch.Buckets)))

				if chi > 5*float64(len(sketch.Buckets)) &&
					(1-h/hMax) > 0.15 {

					fmt.Printf(
						"STRUCTURE w=%d shift=%d chi=%.1f entropy-loss=%.2f\n",
						windowBits,
						shift,
						chi,
						1-h/hMax,
					)

					_ = cmd.PlotFrequencySketch(
						sketch,
						windowBits,
						shift,
						outDir,
					)
				}
			}
		}

	default:
		fmt.Println("Unknown mode:", mode)
	}
}
