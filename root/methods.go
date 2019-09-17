package root

import (
	"math"

	"github.com/fako1024/numerics"
)

const (
	bisectTolerance = 1e-11
	bisectMaxIter   = 100
)

// Linear root finding methods

// Bisect performs a simple bisection of a function within a lower and an upper
// limit
func Bisect(fx func(x float64) float64, aInit, bInit float64) float64 {

	// Define current lower / upper limits based on input parameters
	a, b := aInit, bInit

	// Bisection loop
	for i := 0; i < bisectMaxIter; i++ {

		// Split the current interval in half
		c := (a + b) / 2.

		// Determine the function value at that value, return if within tolerance of
		// expectation
		fxVal := fx(c)
		if fxVal == 0 || (b-a)/2. < bisectTolerance {
			return c
		}

		// Otherwise follow sign of function value vs. sign of expectation
		if numerics.Sign(fxVal) == numerics.Sign(fx(a)) {
			a = c
		} else {
			b = c
		}
	}

	// If bisection failed, return NaN
	return math.NaN()
}

// Non-linear root finding methods

// Method wraps the functional parameters used in root finding methods in a more
// readable type
type Method func(x float64, fx, dfx func(float64) float64) float64

// NewtonRaphson performs the original method by Newton / Raphson
func NewtonRaphson(x float64, fx, dfx func(float64) float64) float64 {
	return x - fx(x)/dfx(x)
}

// Homeier performs a modified Newton method with cubic convergence as introduced
// in "A modified Newton method for rootfinding with cubic convergence", Journal
// of Computational and Applied Mathematics 157 (2003) 227–230
// doi:10.1016/S0377-0427(03)00391-1
func Homeier(x float64, fx, dfx func(float64) float64) float64 {
	fxVal := fx(x)
	return x - fxVal/dfx(x-0.5*fxVal/dfx(x))
}
