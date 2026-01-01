package cmd

import "math"

const (
	ChiFactor      = 5.0  // 5× Freiheitsgrade
	EntropyLossMin = 0.15 // ≥15% Entropieverlust
)

func chiSquare(hist map[uint64]uint64) float64 {
	var total uint64
	for _, v := range hist {
		total += v
	}
	if total == 0 || len(hist) < 2 {
		return 0
	}

	expected := float64(total) / float64(len(hist))
	var chi float64

	for _, v := range hist {
		diff := float64(v) - expected
		chi += (diff * diff) / expected
	}
	return chi
}

func entropy(hist map[uint64]uint64) float64 {
	var total uint64
	for _, v := range hist {
		total += v
	}
	if total == 0 {
		return 0
	}

	var h float64
	for _, v := range hist {
		p := float64(v) / float64(total)
		h -= p * math.Log2(p)
	}
	return h
}

func ChiSquareSketch(s *Sketch) float64 {
	var total uint64
	for _, v := range s.Buckets {
		total += uint64(v)
	}

	expected := float64(total) / float64(len(s.Buckets))
	var chi float64

	for _, v := range s.Buckets {
		d := float64(v) - expected
		chi += (d * d) / expected
	}
	return chi
}

func EntropySketch(s *Sketch) float64 {
	var total float64
	for _, v := range s.Buckets {
		total += float64(v)
	}

	var h float64
	for _, v := range s.Buckets {
		if v == 0 {
			continue
		}
		p := float64(v) / total
		h -= p * math.Log2(p)
	}
	return h
}

func IsClearlyNonUniform(hist map[uint64]uint64) (bool, float64, float64) {
	k := float64(len(hist)) - 1
	if k <= 0 {
		return false, 0, 0
	}

	chi := chiSquare(hist)
	h := entropy(hist)
	hMax := math.Log2(float64(len(hist)))
	entropyLoss := 1 - (h / hMax)

	if chi > ChiFactor*k && entropyLoss > EntropyLossMin {
		return true, chi, entropyLoss
	}
	return false, chi, entropyLoss
}
