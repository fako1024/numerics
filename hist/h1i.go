package hist

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"text/tabwriter"
)

// H1I denotes a one-dimensional histogram
type H1I struct {
	nEntries int
	nBins    int

	sumOfWeights float64

	binContent  []float64
	binVariance []float64
	bins        []float64
}

// NewH1I instantiates a new one-dimensional histogram
func NewH1I(binCenters []float64) *H1I {
	obj := H1I{
		nBins: len(binCenters),

		binContent:  make([]float64, len(binCenters)+2),
		binVariance: make([]float64, len(binCenters)+2),
		bins:        binCenters,
	}

	return &obj
}

// Print prints out the histogram data
func (h *H1I) Print(w io.Writer) error {

	tabw := tabwriter.NewWriter(w, 2, 2, 2, byte(' '), 0)

	yfmt := func(y float64) string {
		if y > 0 {
			return strconv.Itoa(int(y))
		}
		return ""
	}

	fmt.Fprintf(w, "Mode: %.2f\n", h.Mode())

	for i := 0; i < len(h.bins); i++ {
		fmt.Fprintf(tabw, "%s\t%.3g%%\t%s\n",
			fmt.Sprintf("%.4g", h.bins[i]),
			float64(h.BinContent(i+1))*100.0/float64(h.sumOfWeights),
			bar(float64(h.BinContent(i+1))*100.0/float64(h.sumOfWeights))+"\t"+yfmt(h.BinContent(i+1)),
		)
	}

	return tabw.Flush()

}

// NBins Returns the number of bins in the histogram
func (h *H1I) NBins() int {
	return h.nBins
}

// NEntries returns the number of entries in the histogram
func (h *H1I) NEntries() int {
	return h.nEntries
}

// Sum returns the sum of weights in the histogram
func (h *H1I) Sum() float64 {
	return h.sumOfWeights
}

// XMin returns the lower boundary of the x axis
func (h *H1I) XMin() float64 {
	return h.bins[0]
}

// Xmax returns the upper boundary of the x axis
func (h *H1I) XMax() float64 {
	return h.bins[h.nBins-1]
}

// BinContent returns the sum of weights in a particular bin
func (h *H1I) BinContent(bin int) float64 {
	return h.binContent[bin]
}

// BinVariance returns the variance in a particular bin
func (h *H1I) BinVariance(bin int) float64 {
	return h.binVariance[bin]
}

// ModeBin returns the maximum bin
func (h *H1I) MaximumBin() int {
	max, maxBin := -1e99, 0

	for i := 0; i < len(h.bins); i++ {
		if h.binContent[i+1] > max {
			max = h.binContent[i+1]
			maxBin = i
		}
	}

	return maxBin
}

// BinCenter returns the center x value of a particular bin
func (h *H1I) BinCenter(bin int) float64 {
	return h.bins[bin]
}

// Mode returns the mode of the histogram
func (h *H1I) Mode() float64 {
	return h.BinCenter(h.MaximumBin())
}

// SetBinContent sets the sum of weights in a particular bin
func (h *H1I) SetBinContent(bin int, sumOfWeights float64) {

	// increase overall sum of weights by current value in requested bin and
	// subtract the old bin content
	h.sumOfWeights += sumOfWeights - h.binContent[bin]

	h.binContent[bin] = sumOfWeights
}

// SetBinVariance sets the variance in a particular bin
func (h *H1I) SetBinVariance(bin int, variance float64) {
	h.binVariance[bin] = variance
}

// Fill adds a weight / entry to the histogram
func (h *H1I) Fill(val float64, weight ...float64) {

	if len(weight) > 1 {
		panic("must specify no or exactly one weight")
	}
	w := 1.0
	if len(weight) == 1 {
		w = weight[0]
	}

	// Increment counters
	h.nEntries++
	h.sumOfWeights += w

	// Handle underflow case
	if val < h.bins[0] {
		h.binContent[0] += w
		return
	}

	// Handle overflow case
	if val > h.bins[h.nBins-1] {
		h.binContent[h.nBins+1] += w
		return
	}

	// Handle standard case
	for i := 0; i < h.nBins; i++ {
		if almostEqual(val, h.bins[i]) {
			h.binContent[i+1] += w
			return
		}
	}

	panic("invalid value")
}

// Scale scales the histogram by a constant factor
func (h *H1I) Scale(scale float64) {

	h.sumOfWeights *= scale

	for i := 0; i < h.nBins+2; i++ {
		h.binContent[i] *= scale
		h.binVariance[i] *= scale
	}
}

// FindBin returns the bin best matching the value x
func (h *H1I) FindBin(x float64) int {

	if x < h.XMin() {
		return 0
	}
	if x > h.XMax() {
		return h.nBins + 1
	}

	return 1 + int(float64(h.nBins)*(x-h.XMin())/(h.XMax()-h.XMin()))
}

// Interpolate linearly interpolates between the nearest bin neigbors
func (h *H1I) Interpolate(x float64) float64 {

	xBin := h.FindBin(x)

	if x <= h.BinCenter(1) {
		return h.BinContent(1)
	}
	if x >= h.BinCenter(h.NBins()) {
		return h.BinContent(h.NBins())
	}

	var x0, y0, x1, y1 float64
	if x <= h.BinCenter(xBin) {
		y0 = h.BinContent(xBin - 1)
		x0 = h.BinCenter(xBin - 1)
		y1 = h.BinContent(xBin)
		x1 = h.BinCenter(xBin)
	} else {
		y0 = h.BinContent(xBin)
		x0 = h.BinCenter(xBin)
		y1 = h.BinContent(xBin + 1)
		x1 = h.BinCenter(xBin + 1)
	}

	return y0 + (x-x0)*((y1-y0)/(x1-x0))
}

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
