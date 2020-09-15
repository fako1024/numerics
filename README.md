# Method / Tools for numerical methods / statistics

[![Github Release](https://img.shields.io/github/release/fako1024/numerics.svg)](https://github.com/fako1024/numerics/releases)
[![GoDoc](https://godoc.org/github.com/fako1024/numerics?status.svg)](https://godoc.org/github.com/fako1024/numerics/)
[![Go Report Card](https://goreportcard.com/badge/github.com/fako1024/numerics)](https://goreportcard.com/report/github.com/fako1024/numerics)
[![Build/Test Status](https://github.com/fako1024/numerics/workflows/Go/badge.svg)](https://github.com/fako1024/numerics/actions?query=workflow%3AGo)

This package provides several methods / tools that support numerical methods and statistics in Go.

## Features
- Various numeric methods, such as
	- Complete, incomplete and regularized incomplete Beta function
	- Binomial distribution function
	- Sign function
	- Lgamma function (without error return for ease of use)
- Numerical root finding methods (sub-package `root`) via a generic interface, including
	- Linear root finding via Bisection
	- Non-linear root finding via Newton-Raphson and a cubic method
	- Adaptive / heuristic options to circumvent known limitations of root finding methods, i.e. detection of stationary and cyclic situations

## Installation
```bash
go get -u github.com/fako1024/numerics
```

## API summary

The API of the package is fairly straight-forward. The following functions are exposed:
```Go
// Sign returns the sign of a float64  
func Sign(x float64) int

// Lgamma return the logarithmic Gamma function (ignoring any error)  
func Lgamma(x float64) float64

// Beta returns the value of the complete beta function B(a, b).  
func Beta(a, b float64) float64

// BetaIncomplete returns the value of the regularized incomplete beta  
// function Iₓ(a, b).  
//  
// This is not to be confused with the "incomplete beta function",  
// which can be computed as BetaIncomplete(x, a, b)*Beta(a, b).  
//  
// If x < 0 or x > 1, returns NaN.  
func BetaIncompleteRegular(x, a, b float64) float64

// BetaIncomplete returns the value of the (non-regularized) incomplete beta function  
func BetaIncomplete(x, a, b float64) float64

// Binomial returns the value of the probability distribution for a Bernoulli experiment.  
// Consequentially, this is also the differentiated value of the regularized incomplete  
// beta function, representing the cumulative distribution of the binomial PDF  
func Binomial(x, k, n float64) float64
```
The documentation for root finding methods can be found in the sub-package `root`.

## Examples
For some simple examples, have a look at the `numerics_test.go` file.
