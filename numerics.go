package numerics

import (
	"math"
)

// Sign returns the sign of a float64
func Sign(x float64) int {
	if x < 0. {
		return -1
	}
	if x > 0. {
		return 1
	}
	return 0
}

// Lgamma return the logarithmic Gamma function (ignoring any error)
func Lgamma(x float64) float64 {
	y, _ := math.Lgamma(x)
	return y
}

// Beta returns the value of the complete beta function B(a, b).
func Beta(a, b float64) float64 {

	// B(x,y) = Γ(x)Γ(y) / Γ(x+y)
	return math.Exp(Lgamma(a) + Lgamma(b) - Lgamma(a+b))
}

// BetaIncompleteRegular returns the value of the regularized incomplete beta
// function Iₓ(a, b).
//
// This is not to be confused with the "incomplete beta function",
// which can be computed as BetaIncomplete(x, a, b)*Beta(a, b).
//
// If x < 0 or x > 1, returns NaN.
func BetaIncompleteRegular(x, a, b float64) float64 {

	// Based on Numerical Recipes in C, section 6.4. This uses the
	// continued fraction definition of I:
	//
	//  (xᵃ*(1-x)ᵇ)/(a*B(a,b)) * (1/(1+(d₁/(1+(d₂/(1+...))))))
	//
	// where B(a,b) is the beta function and
	//
	//  d_{2m+1} = -(a+m)(a+b+m)x/((a+2m)(a+2m+1))
	//  d_{2m}   = m(b-m)x/((a+2m-1)(a+2m))
	if x < 0 || x > 1 {
		return math.NaN()
	}
	bt := 0.0
	if 0 < x && x < 1 {

		// Compute the coefficient before the continued
		// fraction.
		bt = math.Exp(Lgamma(a+b) - Lgamma(a) - Lgamma(b) +
			a*math.Log(x) + b*math.Log(1-x))
	}
	if x < (a+1)/(a+b+2) {
		// Compute continued fraction directly.
		return bt * betacf(x, a, b) / a
	}

	// Compute continued fraction after symmetry transform.
	return 1 - bt*betacf(1-x, b, a)/b
}

// BetaIncomplete returns the value of the (non-regularized) incomplete beta function
func BetaIncomplete(x, a, b float64) float64 {
	return BetaIncompleteRegular(x, a, b) * Beta(a, b)
}

// Binomial returns the value of the probability distribution for a Bernoulli experiment.
// Consequentially, this is also the differentiated value of the regularized incomplete
// beta function, representing the cumulative distribution of the binomial PDF
func Binomial(x, k, n float64) float64 {

	// ∂ₓIₓ(k, n) = exp( (n-k) * ln(1-x) + k * ln(x) )
	return math.Exp((n-k)*math.Log(1.-x) + k*math.Log(x))
}

////////////////////////////////////////////////////////////////////////////////

const (
	betaEpsilon       = 3e-14
	betaMaxIterations = 200
)

// smallestNonZero return the smalles non-zero value to avoid creating division
// by zero situations due to numeric fluctuations
func smallestNonZero(val float64) float64 {
	if math.Abs(val) < math.SmallestNonzeroFloat64 {
		return math.SmallestNonzeroFloat64
	}

	return val
}

// betacf is the continued fraction component of the regularized
// incomplete beta function Iₓ(a, b).
// Based on Numerical Recipes in C, Second Edition, Section 6.4
func betacf(x, a, b float64) float64 {

	c := 1.0
	d := 1.0 / smallestNonZero(1.0-(a+b)*x/(a+1.0))
	h := d
	for m := 1; m <= betaMaxIterations; m++ {
		mf := float64(m)

		// One step (the even one) of the recurrence
		numer := mf * (b - mf) * x / ((a + 2.0*mf - 1.0) * (a + 2.0*mf))
		d = 1 / smallestNonZero(1+numer*d)
		c = smallestNonZero(1 + numer/c)
		h *= d * c

		// Next step of the recurrence (the odd one)
		numer = -(a + mf) * (a + b + mf) * x / ((a + 2*mf) * (a + 2.0*mf + 1.0))
		d = 1 / smallestNonZero(1+numer*d)
		c = smallestNonZero(1 + numer/c)
		hfac := d * c
		h *= hfac

		// If sufficient precision is reached, return
		if math.Abs(hfac-1) < betaEpsilon {
			return h
		}
	}

	// If function did not converge, return NaN
	return math.NaN()
}
