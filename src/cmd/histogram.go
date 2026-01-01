package cmd

import (
	"os"
)

type Endian int

type Sketch struct {
	Buckets []uint32
	Mask    uint64
}

const (
	BigEndian Endian = iota
	LittleEndian
)

func LoadBits(filename string, endian Endian) []uint8 {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	bits := make([]uint8, 0, len(data)*8)

	for _, b := range data {
		if endian == BigEndian {
			for i := 7; i >= 0; i-- {
				bits = append(bits, (b>>i)&1)
			}
		} else {
			for i := 0; i < 8; i++ {
				bits = append(bits, (b>>i)&1)
			}
		}
	}
	return bits
}

func NewSketch(bits int) *Sketch {
	size := 1 << bits
	return &Sketch{
		Buckets: make([]uint32, size),
		Mask:    uint64(size - 1),
	}
}

func FillSketch(bits []uint8, windowBits int, shift int, s *Sketch) {
	if shift+windowBits >= len(bits) {
		return
	}

	var w uint64
	mask := uint64((1 << windowBits) - 1)

	for i := 0; i < windowBits; i++ {
		w = (w << 1) | uint64(bits[shift+i])
	}
	s.Buckets[w&s.Mask]++

	for i := shift + windowBits; i < len(bits); i++ {
		w = ((w << 1) & mask) | uint64(bits[i])
		s.Buckets[w&s.Mask]++
	}
}

func GetHistogram(filename string, endian Endian) [8][256]uint64 {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var histograms [8][256]uint64

	// Convert entire file into a bit slice first
	var bits []uint8

	for _, b := range data {
		if endian == BigEndian {
			for i := 7; i >= 0; i-- {
				bits = append(bits, (b>>i)&1)
			}
		} else {
			for i := 0; i < 8; i++ {
				bits = append(bits, (b>>i)&1)
			}
		}
	}

	// For each shift
	for shift := 0; shift < 8; shift++ {
		for i := shift; i+7 < len(bits); i += 8 {
			var window uint8 = 0
			for j := 0; j < 8; j++ {
				window = (window << 1) | bits[i+j]
			}
			histograms[shift][window]++
		}
	}

	return histograms
}

func GetHistogramWindows(
	filename string,
	endian Endian,
	windowBits int,
) []map[uint64]uint64 {

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var bits []uint8
	for _, b := range data {
		if endian == BigEndian {
			for i := 7; i >= 0; i-- {
				bits = append(bits, (b>>i)&1)
			}
		} else {
			for i := 0; i < 8; i++ {
				bits = append(bits, (b>>i)&1)
			}
		}
	}

	hists := make([]map[uint64]uint64, windowBits)

	for shift := 0; shift < windowBits; shift++ {
		hist := make(map[uint64]uint64)

		for i := shift; i+windowBits <= len(bits); i += windowBits {
			var w uint64
			for j := 0; j < windowBits; j++ {
				w = (w << 1) | uint64(bits[i+j])
			}
			hist[w]++
		}
		hists[shift] = hist
	}
	return hists
}
