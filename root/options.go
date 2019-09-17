package root

// WithMinIterations sets a minimum number of iterations to perform
func WithMinIterations(nIterations int) func(*Finder) {
	return func(n *Finder) {
		n.minIterations = nIterations
	}
}

// WithMaxIterations sets a maximum number of iterations to perform
func WithMaxIterations(nIterations int) func(*Finder) {
	return func(n *Finder) {
		n.maxIterations = nIterations
	}
}

// WithTargetPrecision sets a target precision (max. deviation from target x) for
// the method, implicitly determining the number of iterations to be performed
func WithTargetPrecision(targetPrecision float64) func(*Finder) {
	return func(n *Finder) {
		n.targetPrecision = targetPrecision
	}
}

// WithMethod sets a specific method to be used to perform the iterative process
func WithMethod(method Method) func(*Finder) {
	return func(n *Finder) {
		n.method = method
	}
}

// WithLimits sets limits for the variable x
func WithLimits(xMin, xMax float64) func(*Finder) {
	return func(n *Finder) {
		n.xMin, n.xMax = xMin, xMax
	}
}

// WithHeuristics enables adaptive methods to circumvent known limitations of the
// Newton-Raphson method, i.e. detection of stationary and cyclic situations
func WithHeuristics() func(*Finder) {
	return func(n *Finder) {
		n.useHeuristics = true
	}
}
