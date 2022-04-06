package promqlext

import (
	"math"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
)

const okZeroIfNone = "ok_zero_if_none"

// RegisterOKZeroIfNone регистрирует функцию `ok_zero_if_none` в promql
//
// Функция `ok_zero_if_none` это оконная функция, зануляющая значение компонент вектора, данные
// которых устарели (старше одной минуты ActualInterval).
//
// Вторым параметром передаётся время (в виде мс), раньше которого данные игнорируются.
func RegisterOKZeroIfNone() {
	parser.Functions[okZeroIfNone] = &parser.Function{
		Name:       okZeroIfNone,
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix, parser.ValueTypeScalar},
		Variadic:   1,
		ReturnType: parser.ValueTypeVector,
	}
	promql.FunctionCalls[okZeroIfNone] = funcOKZeroIfNone
}

// === ok_zero_if_none(Matrix parser.ValueTypeMatrix, Ms parser.ValueTypeScalar) Vector ===
func funcOKZeroIfNone(vals []parser.Value, args parser.Expressions, enh *promql.EvalNodeHelper) promql.Vector {
	vec := vals[0].(promql.Matrix)
	var ms int64 = math.MinInt64
	if len(vals) > 1 {
		ms = int64(vals[1].(promql.Vector)[0].V)
	}
	for _, el := range vec {
		p := takeLast(el)
		if p.T < ms {
			continue
		}
		value := p.V
		if isNone(p, enh) || math.IsNaN(value) {
			value = 0
		}
		enh.Out = append(enh.Out, promql.Sample{
			Metric: enh.DropMetricName(el.Metric),
			Point:  promql.Point{V: value},
		})
	}
	return enh.Out
}
