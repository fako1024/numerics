
# Numerical methods for function root finding
This package provides numerical methods for finding the root of functions using linear and quadratic / cubic approaches.

## Features
- Numerical root finding methods via a generic interface, including
	- Linear root finding via Bisection
	- Non-linear root finding via Newton-Raphson and a cubic method
	- Adaptive / heuristic options to circumvent known limitations of root finding methods, i.e. detection of stationary and cyclic situations

## Installation
```bash
go get -u github.com/fako1024/numerics/root
```

## API summary

The API of the package is fairly straight-forward. The following functions are exposed:
```Go
// Bisect performs a simple bisection of a function within a lower and an upper
// limit
func Bisect(fx func(x float64) float64, aInit, bInit float64) float64

/////////////////

// Method wraps the functional parameters used in root finding methods in a more
// readable type
type Method func(x float64, fx, dfx func(float64) float64) float64

// Find perform a non-linear iterative root-finding method using the
// provided parameters / options
func Find(fx, dfx func(x float64) float64, xInit float64, options ...func(*Finder)) float64

/////////////////

// NewtonRaphson performs the original method by Newton / Raphson
func NewtonRaphson(x float64, fx, dfx func(float64) float64) float64

// Homeier performs a modified Newton method with cubic convergence as introduced
// in "A modified Newton method for rootfinding with cubic convergence", Journal
// of Computational and Applied Mathematics 157 (2003) 227–230
// doi:10.1016/S0377-0427(03)00391-1
func Homeier(x float64, fx, dfx func(float64) float64) float64
```

## Examples
For some simple examples, have a look at the `root_test.go` file.
