package promqlext

import (
	"math"

	"github.com/prometheus/prometheus/model/value"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
)

const okSmoothie = "ok_smoothie"

// RegisterOKDefined регистрирует функцию `ok_smoothie` в promql
//
// Функция `ok_smoothie` это оконная функция, заменяющая значение компонент вектора
// на среднее за период.
func RegisterOKSmoothie() {
	parser.Functions[okSmoothie] = &parser.Function{
		Name:       okSmoothie,
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
		ReturnType: parser.ValueTypeVector,
	}
	promql.FunctionCalls[okSmoothie] = funcOKSmoothie
}

func aggrOverTime(vals []parser.Value, enh *promql.EvalNodeHelper, aggrFn func([]promql.Point) float64) promql.Vector {
	el := vals[0].(promql.Matrix)[0]
	v := aggrFn(el.Points)
	if value.IsStaleNaN(v) {
		return enh.Out
	}

	return append(enh.Out, promql.Sample{
		Point: promql.Point{V: v},
	})
}

// === ok_smoothie(Matrix parser.ValueTypeMatrix) Vector ===
func funcOKSmoothie(vals []parser.Value, args parser.Expressions, enh *promql.EvalNodeHelper) promql.Vector {
	return aggrOverTime(vals, enh, func(values []promql.Point) float64 {
		var mean, count, c float64
		if len(values) == 0 {
			return math.NaN()
		}
		v := values[len(values)-1]
		if math.IsNaN(v.V) {
			return v.V
		}
		if enh.Ts != v.T {
			return math.Float64frombits(value.StaleNaN)
		}
		for _, v := range values {
			count++
			if math.IsInf(mean, 0) {
				if math.IsInf(v.V, 0) && (mean > 0) == (v.V > 0) {
					// The `mean` and `v.V` values are `Inf` of the same sign.  They
					// can't be subtracted, but the value of `mean` is correct
					// already.
					continue
				}
				if !math.IsInf(v.V, 0) && !math.IsNaN(v.V) {
					// At this stage, the mean is an infinite. If the added
					// value is neither an Inf or a Nan, we can keep that mean
					// value.
					// This is required because our calculation below removes
					// the mean value, which would look like Inf += x - Inf and
					// end up as a NaN.
					continue
				}
			}
			mean, c = kahanSumInc(v.V/count-mean/count, mean, c)
		}

		if math.IsInf(mean, 0) {
			return mean
		}
		return mean + c
	})
}

func kahanSumInc(inc, sum, c float64) (newSum, newC float64) {
	t := sum + inc
	// Using Neumaier improvement, swap if next term larger than sum.
	if math.Abs(sum) >= math.Abs(inc) {
		c += (sum - t) + inc
	} else {
		c += (inc - t) + sum
	}
	return t, c
}
