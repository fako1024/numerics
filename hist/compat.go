package hist

// H1D denotes a one-dimensional histogram based on float64 values
type H1D = H1[float64]

// NewH1D instantiates a new one-dimensional histogram based on float64 values
func NewH1D(n int, xMin, xMax float64) *H1D {
	return NewH1(n, xMin, xMax)
}

// H1I denotes a one-dimensional histogram based on integer values
type H1I = H1[int]

// NewH1I instantiates a new one-dimensional histogram based on integer values
func NewH1I(n int, xMin, xMax int) *H1I {
	return NewH1(n, xMin, xMax)
}
