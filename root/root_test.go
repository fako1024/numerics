package root

import (
	"fmt"
	"math"
	"path"
	"reflect"
	"runtime"
	"testing"
)

const expectedPrecision = 1e-9

var methods = []Method{
	NewtonRaphson,
	Homeier,
}

type testCaseBisect struct {
	fx         func(float64) float64
	xMin, xMax float64
}

type testCaseNewton struct {
	fx, dfx func(float64) float64
	xInit   float64
}

func TestOptions(t *testing.T) {
	_ = Find(func(x float64) float64 {
		return x*x - 612
	}, func(x float64) float64 {
		return 2 * x
	}, 10.,
		WithHeuristics(),
		WithMethod(Homeier),
		WithMinIterations(5),
		WithMaxIterations(25),
		WithTargetPrecision(1e-9),
		WithLimits(-1e9, 1e9),
	)
}

func TestBisectTable(t *testing.T) {

	testCases := map[string]testCaseBisect{
		"SquareRoot2": {
			fx: func(x float64) float64 {
				return x*x - 612
			},
			xMin: 1.,
			xMax: 50.,
		},
		"CosineEquation": {
			fx: func(x float64) float64 {
				return math.Cos(x) - x*x*x
			},
			xMin: 0.1,
			xMax: 1.0,
		},
	}

	for testName, cs := range testCases {
		t.Run(testName, func(t *testing.T) {
			root := Bisect(cs.fx, cs.xMin, cs.xMax)

			if math.IsNaN(root) || math.IsInf(root, 0) {
				t.Fatalf("Unexpected non-numerical result for %s: %v", testName, root)
			}

			if math.Abs(cs.fx(root)) > expectedPrecision {
				t.Fatalf("Estimated value of f(x) for %s deviates significantly from expectation: have %.5f, want 0", testName, cs.fx(root))
			}
		})
	}

}

func TestNewtonTable(t *testing.T) {

	testCases := map[string]testCaseNewton{
		"SquareRoot2": {
			fx: func(x float64) float64 {
				return x*x - 612
			},
			dfx: func(x float64) float64 {
				return 2 * x
			},
			xInit: 10.,
		},
		"CosineEquation": {
			fx: func(x float64) float64 {
				return math.Cos(x) - x*x*x
			},
			dfx: func(x float64) float64 {
				return -math.Sin(x) - 3*x*x
			},
			xInit: 0.5,
		},
	}

	for testName, cs := range testCases {
		for _, method := range methods {
			t.Run(caseName(method, testName), func(t *testing.T) {
				root := Find(cs.fx, cs.dfx, cs.xInit, WithHeuristics(), WithMethod(method))

				if math.IsNaN(root) || math.IsInf(root, 0) {
					t.Fatalf("Unexpected non-numerical result for %s: %v", testName, root)
				}

				if math.Abs(cs.fx(root)) > expectedPrecision {
					t.Fatalf("Estimated value of f(x) for %s deviates significantly from expectation: have %.5f, want 0", testName, cs.fx(root))
				}
			})
		}
	}
}

func TestTableHeuristics(t *testing.T) {

	testCases := map[string]testCaseNewton{
		"StationaryPoint": {
			fx: func(x float64) float64 {
				return 1. - x*x
			},
			dfx: func(x float64) float64 {
				return -2. * x
			},
			xInit: 0.,
		},
		"Cycle": {
			fx: func(x float64) float64 {
				return x*x*x - 2.*x + 2.
			},
			dfx: func(x float64) float64 {
				return 3*x*x - 2.
			},
			xInit: 0.,
		},
	}

	for testName, cs := range testCases {
		for _, method := range methods {
			t.Run(caseName(method, testName), func(t *testing.T) {
				root := Find(cs.fx, cs.dfx, cs.xInit, WithHeuristics())

				if math.IsNaN(root) || math.IsInf(root, 0) {
					t.Fatalf("Unexpected non-numerical result for %s: %v", testName, root)
				}

				if math.Abs(cs.fx(root)) > expectedPrecision {
					t.Fatalf("Estimated value of f(x) for %s deviates significantly from expectation: have %.5f, want 0", testName, cs.fx(root))
				}
			})
		}
	}
}

func BenchmarkMethods(b *testing.B) {
	cs := testCaseNewton{
		fx: func(x float64) float64 {
			return math.Cos(x) - x*x*x
		},
		dfx: func(x float64) float64 {
			return -math.Sin(x) - 3*x*x
		},
		xInit: 0.5,
	}

	for _, method := range methods {
		b.Run(caseName(method, "bench"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				root := Find(cs.fx, cs.dfx, cs.xInit, WithMethod(method), WithMinIterations(2))
				_ = root
			}
		})
	}
}

func caseName(i interface{}, suffix string) string {
	return fmt.Sprintf("%s_%s", path.Base(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()), suffix)
}
