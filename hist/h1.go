package hist

import (
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"
	"time"
)

// Number provides a type constraint on the supported generics (anything number-like)
type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | time.Duration | uintptr
}

// H1 denotes a one-dimensional histogram
type H1[T Number] struct {
	nEntries int
	nBins    int

	sumOfWeights float64

	binContent  []float64
	binVariance []float64
	bins        []T
}

// NewH1 instantiates a new one-dimensional histogram
func NewH1[T Number](n int, xMin, xMax T) *H1[T] {
	obj := H1[T]{
		nBins: n,

		binContent:  make([]float64, n+2),
		binVariance: make([]float64, n+2),
		bins:        make([]T, n+1),
	}

	step := (xMax - xMin) / T(n)
	for i := 0; i < n+1; i++ {
		obj.bins[i] = xMin + T(i)*step
	}

	return &obj
}

// Print prints out the histogram data to any io.Writer
func (h *H1[T]) Print(w io.Writer) error {

	tabw := tabwriter.NewWriter(w, 2, 2, 2, byte(' '), 0)

	yfmt := func(y float64) string {
		if y > 0 {
			return strconv.Itoa(int(y))
		}
		return ""
	}

	fmt.Fprintf(w, "Mode: %v\n", h.Mode())

	for i := 0; i < len(h.bins)-1; i++ {
		fmt.Fprintf(tabw, "%s-%s\t%.3g%%\t%s\n",
			fmt.Sprintf("%.4v", h.bins[i]),
			fmt.Sprintf("%.4v", h.bins[i+1]),
			h.BinContent(i+1)*100.0/h.sumOfWeights,
			bar(h.BinContent(i+1)*100.0/h.sumOfWeights)+"\t"+yfmt(h.BinContent(i+1)),
		)
	}

	return tabw.Flush()

}

// NBins Returns the number of bins in the histogram
func (h *H1[T]) NBins() int {
	return h.nBins
}

// NEntries returns the number of entries in the histogram
func (h *H1[T]) NEntries() int {
	return h.nEntries
}

// Sum returns the sum of weights in the histogram
func (h *H1[T]) Sum() float64 {
	return h.sumOfWeights
}

// XMin returns the lower boundary of the x axis
func (h *H1[T]) XMin() T {
	return h.bins[0]
}

// XMax returns the upper boundary of the x axis
func (h *H1[T]) XMax() T {
	return h.bins[h.nBins]
}

// BinContent returns the sum of weights in a particular bin
func (h *H1[T]) BinContent(bin int) float64 {
	return h.binContent[bin]
}

// BinVariance returns the variance in a particular bin
func (h *H1[T]) BinVariance(bin int) float64 {
	return h.binVariance[bin]
}

// MaximumBin returns the maximum bin
func (h *H1[T]) MaximumBin() int {
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
func (h *H1[T]) BinCenter(bin int) T {
	return (h.bins[bin-1] + h.bins[bin]) / 2.0
}

// Mode returns the mode of the histogram
func (h *H1[T]) Mode() T {
	return h.BinCenter(h.MaximumBin())
}

// SetBinContent sets the sum of weights in a particular bin
func (h *H1[T]) SetBinContent(bin int, sumOfWeights float64) {

	// increase overall sum of weights by current value in requested bin and
	// subtract the old bin content
	h.sumOfWeights += sumOfWeights - h.binContent[bin]

	h.binContent[bin] = sumOfWeights
}

// SetBinVariance sets the variance in a particular bin
func (h *H1[T]) SetBinVariance(bin int, variance float64) {
	h.binVariance[bin] = variance
}

// Fill adds a weight / entry to the histogram
func (h *H1[T]) Fill(val T, weight ...float64) {

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
func (h *H1[T]) Scale(scale float64) {

	h.sumOfWeights *= scale

	for i := 0; i < h.nBins+2; i++ {
		h.binContent[i] *= scale
		h.binVariance[i] *= scale
	}
}

// FindBin returns the bin best matching the value x
func (h *H1[T]) FindBin(x T) int {

	if x < h.XMin() {
		return 0
	}
	if x > h.XMax() {
		return h.nBins + 1
	}

	return 1 + int(T(h.nBins)*(x-h.XMin())/(h.XMax()-h.XMin()))
}

// Interpolate linearly interpolates between the nearest bin neigbors
func (h *H1[T]) Interpolate(x T) float64 {

	xBin := h.FindBin(x)

	if x <= h.BinCenter(1) {
		return h.BinContent(1)
	}
	if x >= h.BinCenter(h.NBins()) {
		return h.BinContent(h.NBins())
	}

	var (
		x0, x1 T
		y0, y1 float64
	)
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

	return y0 + float64(x-x0)*((y1-y0)/float64(x1-x0))
}
