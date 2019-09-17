package root

import (
	"math"
)

const maxRetries = 100

// Finder defines a non-linear approach to root finding
type Finder struct {
	fx, dfx func(x float64) float64
	method  Method

	xMin, xMax float64

	minIterations   int
	maxIterations   int
	targetPrecision float64
	useHeuristics   bool
}

// Find perform a non-linear iterative root-finding method using the
// provided parameters / options
func Find(fx, dfx func(x float64) float64, xInit float64, options ...func(*Finder)) float64 {

	obj := &Finder{
		fx:     fx,
		dfx:    dfx,
		method: NewtonRaphson,

		xMin: -math.MaxFloat64,
		xMax: math.MaxFloat64,

		minIterations:   5,
		maxIterations:   25,
		targetPrecision: 1e-9,
	}

	// Execute functional options (if any), see options.go for implementation
	for _, option := range options {
		option(obj)
	}

	return obj.loop(xInit)
}

////////////////////////////////////////////////////////////////////////////////

// loop executed the actual root finding loop
func (n *Finder) loop(xInit float64) float64 {

	// Initialize loop variables
	x := xInit
	nIter := 0
	resultLookup := make(map[float64]struct{})

	for {

		// Determine new value for x according to the defined root-finding method
		xNew := n.method(x, n.fx, n.dfx)

		// Guard against excess situations, retrying with a smaller change
		if !math.IsInf(xNew, 0) {
			if xNew > n.xMax {

				// Upper Excess, setting x to (x + xMax)/2
				x = 0.5 * (x + n.xMax)
				continue
			} else if xNew < n.xMin {

				// Lower Excess, setting x to (x + xMin)/2
				x = 0.5 * (x + n.xMin)
				continue
			}
		}

		// If the current value is NaN, return it
		if math.IsNaN(xNew) {
			return math.NaN()
		}

		// If enabled, perform heuristic approach to circumvent known limitations of the
		// Newton-Raphson method, i.e. detection of stationary and cyclic situations
		if n.useHeuristics {

			// Attempt to recover from infinity situations by adapting the value more slowly
			if math.IsInf(xNew, 0) {
				if math.IsInf(xNew, 1) {
					x += 0.1*x + 0.1
				} else {
					x -= 0.1*x - 0.1
				}

				continue
			}

			// Avoid recurring situations / getting "stuck" by storing values already seen
			// and slightly fluctuating the value if values reaccur
			if math.Abs(xNew-x) > 1e-15 {
				if _, alreadySeen := resultLookup[xNew]; alreadySeen {
					if xNew != x {
						x = (xNew + x) / 2.
					} else {
						x += 0.1*x + 0.1
					}
					continue
				}

				// Store value for later lookups
				resultLookup[xNew] = struct{}{}
			}
		}

		x = xNew
		nIter++

		// If the minimum number of iterations has been performed...
		if nIter >= n.minIterations {

			// ... and target precision has been reached or the maximum number of iterations
			// has been performed, break
			if math.Abs(n.fx(x)) < n.targetPrecision || nIter >= n.maxIterations {
				break
			}
		}
	}

	// Return value from latest successful iteration
	return x
}
