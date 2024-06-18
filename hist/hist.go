package hist

import (
	"io"
	"math"
	"strings"
)

type Hist1D interface {
	Print(w io.Writer) error

	// NBins Returns the number of bins in the histogram
	NBins() int

	// NEntries returns the number of entries in the histogram
	NEntries() int

	// Sum returns the sum of weights in the histogram
	Sum() float64

	// XMin returns the lower boundary of the x axis
	XMin() float64

	// Xmax returns the upper boundary of the x axis
	XMax() float64

	// BinContent returns the sum of weights in a particular bin
	BinContent(bin int) float64

	// BinVariance returns the variance in a particular bin
	BinVariance(bin int) float64

	// MaximumBin returns the maximum bin
	MaximumBin() int

	// BinCenter returns the center x value of a particular bin
	BinCenter(bin int) float64

	// Mode returns the mode of the histogram
	Mode() float64

	// SetBinContent sets the sum of weights in a particular bin
	SetBinContent(bin int, sumOfWeights float64)

	// SetBinVariance sets the variance in a particular bin
	SetBinVariance(bin int, variance float64)

	// Fill adds a weight / entry to the histogram
	Fill(val float64, weight ...float64)

	// Scale scales the histogram by a constant factor
	Scale(scale float64)

	// FindBin returns the bin best matching the value x
	FindBin(x float64) int

	// Interpolate linearly interpolates between the nearest bin neigbors
	Interpolate(x float64) float64
}

////////////////////////////////////////////////////////////////////////////////////////////

var blocks = []string{
	"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█",
}

func bar(v float64) string {
	if v < 0. || math.IsNaN(v) {
		v = 0.
	}

	charIdx := int(math.Floor((v-math.Floor(v))*10.0) / 10.0 * 8.0)
	return strings.Repeat("█", int(v)) + blocks[charIdx]
}
