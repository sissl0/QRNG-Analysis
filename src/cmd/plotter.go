package cmd

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func PlotFrequency(freq [256]uint64, out string) error {
	p := plot.New()
	p.Title.Text = "Character Frequency Distribution"
	p.X.Label.Text = "Byte (0–255)"
	p.Y.Label.Text = "Frequency (log)"

	values := make(plotter.Values, 256)
	for i, v := range freq {
		if v == 0 {
			values[i] = math.NaN()
			continue
		}
		values[i] = float64(v)
	}

	bars, err := plotter.NewBarChart(values, vg.Points(2))
	if err != nil {
		return err
	}
	bars.Color = color.RGBA{B: 255, A: 255}

	p.Add(bars)
	p.X.Min = -1
	p.X.Max = 256

	return p.Save(14*vg.Centimeter, 6*vg.Centimeter, out)
}

func PlotFrequencyMap(hist map[uint64]uint64, out string) error {
	p := plot.New()
	p.Title.Text = "Window Frequency Distribution"
	p.X.Label.Text = "Symbol"
	p.Y.Label.Text = "Frequency"

	values := make(plotter.Values, 0, len(hist))
	for _, v := range hist {
		values = append(values, float64(v))
	}

	bars, err := plotter.NewBarChart(values, vg.Points(1))
	if err != nil {
		return err
	}

	p.Add(bars)
	return p.Save(16*vg.Centimeter, 6*vg.Centimeter, out)
}

func PlotFrequencySketch(
	s *Sketch,
	windowBits int,
	shift int,
	outDir string,
) error {

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = fmt.Sprintf("Sketch Distribution w=%d shift=%d", windowBits, shift)
	p.X.Label.Text = "Bucket"
	p.Y.Label.Text = "Frequency"

	values := make(plotter.Values, 0, len(s.Buckets))
	for _, v := range s.Buckets {
		if v > 0 {
			values = append(values, float64(v))
		}
	}
	sort.Float64s(values)

	bars, err := plotter.NewBarChart(values, vg.Points(0.5))
	if err != nil {
		return err
	}
	bars.Color = color.RGBA{B: 255, A: 255}

	p.Add(bars)

	outfile := filepath.Join(
		outDir,
		fmt.Sprintf("cand_w%d_shift%d.png", windowBits, shift),
	)

	return p.Save(18*vg.Centimeter, 6*vg.Centimeter, outfile)
}
