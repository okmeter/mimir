package promqlext

import (
	"math"

	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
)

const okReplaceNaN = "ok_replace_nan"

// RegisterOKReplaceNaN регистрирует функцию `ok_replace_nan` в promql
//
// Функция `ok_replace_nan` это оконная функция, заменяющая значение компонент вектора, данные
// которых устарели (старше одной минуты ActualInterval) на значение, переданное во втором параметре.
//
// Третьим параметром передаётся время (в виде мс), раньше которого данные игнорируются.
func RegisterOKReplaceNaN() {
	parser.Functions[okReplaceNaN] = &parser.Function{
		Name:       okReplaceNaN,
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix, parser.ValueTypeScalar, parser.ValueTypeScalar},
		Variadic:   1,
		ReturnType: parser.ValueTypeVector,
	}
	promql.FunctionCalls[okReplaceNaN] = funcOKReplaceNaN
}

// === ok_replace_nan(Matrix parser.ValueTypeMatrix, Value parser.ValueTypeScalar, Ms parser.ValueTypeScalar) Vector ===
func funcOKReplaceNaN(vals []parser.Value, args parser.Expressions, enh *promql.EvalNodeHelper) promql.Vector {
	vec := vals[0].(promql.Matrix)
	val := vals[1].(promql.Vector)[0].V
	var ms int64 = math.MinInt64
	if len(vals) > 2 {
		ms = int64(vals[2].(promql.Vector)[0].V)
	}
	for _, el := range vec {
		p := takeLast(el)
		if p.T < ms {
			continue
		}
		value := p.V
		if isNone(p, enh) || math.IsNaN(value) {
			value = val
		}
		enh.Out = append(enh.Out, promql.Sample{
			Metric: enh.DropMetricName(el.Metric),
			Point:  promql.Point{V: value},
		})
	}
	return enh.Out
}
