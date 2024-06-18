package hist

import (
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"
)

// H1D denotes a one-dimensional histogram
type H1D struct {
	nEntries int
	nBins    int

	sumOfWeights float64

	binContent  []float64
	binVariance []float64
	bins        []float64
}

// NewH1D instantiates a new one-dimensional histogram
func NewH1D(n int, xMin, xMax float64) *H1D {
	obj := H1D{
		nBins: n,

		binContent:  make([]float64, n+2),
		binVariance: make([]float64, n+2),
		bins:        make([]float64, n+1),
	}

	step := (xMax - xMin) / float64(n)
	for i := 0; i < n+1; i++ {
		obj.bins[i] = xMin + float64(i)*step
	}

	return &obj
}

// Print prints out the histogram data
func (h *H1D) Print(w io.Writer) error {

	tabw := tabwriter.NewWriter(w, 2, 2, 2, byte(' '), 0)

	yfmt := func(y float64) string {
		if y > 0 {
			return strconv.Itoa(int(y))
		}
		return ""
	}

	fmt.Fprintf(w, "Mode: %.2f\n", h.Mode())

	for i := 0; i < len(h.bins)-1; i++ {
		fmt.Fprintf(tabw, "%s-%s\t%.3g%%\t%s\n",
			fmt.Sprintf("%.4g", h.bins[i]),
			fmt.Sprintf("%.4g", h.bins[i+1]),
			float64(h.BinContent(i+1))*100.0/float64(h.sumOfWeights),
			bar(float64(h.BinContent(i+1))*100.0/float64(h.sumOfWeights))+"\t"+yfmt(h.BinContent(i+1)),
		)
	}

	return tabw.Flush()

}

// NBins Returns the number of bins in the histogram
func (h *H1D) NBins() int {
	return h.nBins
}

// NEntries returns the number of entries in the histogram
func (h *H1D) NEntries() int {
	return h.nEntries
}

// Sum returns the sum of weights in the histogram
func (h *H1D) Sum() float64 {
	return h.sumOfWeights
}

// XMin returns the lower boundary of the x axis
func (h *H1D) XMin() float64 {
	return h.bins[0]
}

// Xmax returns the upper boundary of the x axis
func (h *H1D) XMax() float64 {
	return h.bins[h.nBins]
}

// BinContent returns the sum of weights in a particular bin
func (h *H1D) BinContent(bin int) float64 {
	return h.binContent[bin]
}

// BinVariance returns the variance in a particular bin
func (h *H1D) BinVariance(bin int) float64 {
	return h.binVariance[bin]
}

// ModeBin returns the maximum bin
func (h *H1D) MaximumBin() int {
	max, maxBin := -1e99, 0

	for i := 0; i < len(h.bins)-1; i++ {
		if h.binContent[i+1] > max {
			max = h.binContent[i+1]
			maxBin = i + 1
		}
	}

	return maxBin
}

// BinCenter returns the center x value of a particular bin
func (h *H1D) BinCenter(bin int) float64 {
	return (h.bins[bin-1] + h.bins[bin]) / 2.0
}

// Mode returns the mode of the histogram
func (h *H1D) Mode() float64 {
	return h.BinCenter(h.MaximumBin())
}

// SetBinContent sets the sum of weights in a particular bin
func (h *H1D) SetBinContent(bin int, sumOfWeights float64) {

	// increase overall sum of weights by current value in requested bin and
	// subtract the old bin content
	h.sumOfWeights += sumOfWeights - h.binContent[bin]

	h.binContent[bin] = sumOfWeights
}

// SetBinVariance sets the variance in a particular bin
func (h *H1D) SetBinVariance(bin int, variance float64) {
	h.binVariance[bin] = variance
}

// Fill adds a weight / entry to the histogram
func (h *H1D) Fill(val float64, weight ...float64) {

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
	if val > h.bins[h.nBins] {
		h.binContent[h.nBins+1] += w
		return
	}

	// Handle standard case
	for i := 0; i < h.nBins-1; i++ {
		if val >= h.bins[i] && val < h.bins[i+1] {
			h.binContent[i+1] += w
			return
		}
	}

	// Last regular bin is inclusive
	if val >= h.bins[h.nBins-1] && val <= h.bins[h.nBins] {
		h.binContent[h.nBins] += w
	}
}

// Scale scales the histogram by a constant factor
func (h *H1D) Scale(scale float64) {

	h.sumOfWeights *= scale

	for i := 0; i < h.nBins+2; i++ {
		h.binContent[i] *= scale
		h.binVariance[i] *= scale
	}
}

// FindBin returns the bin best matching the value x
func (h *H1D) FindBin(x float64) int {

	if x < h.XMin() {
		return 0
	}
	if x > h.XMax() {
		return h.nBins + 1
	}

	return 1 + h.nBins*int((x-h.XMin())/(h.XMax()-h.XMin()))
}

// Interpolate linearly interpolates between the nearest bin neigbors
func (h *H1D) Interpolate(x float64) float64 {

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
