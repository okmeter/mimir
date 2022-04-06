package promqlext

import (
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
)

const okDefined = "ok_defined"

// RegisterOKDefined регистрирует функцию `ok_defined` в promql
//
// Функция `ok_defined` это оконная функция, заменяющая значение компонент вектора
// на 0, если данные устарели (старше одной минуты ActualInterval) или 1 для актуальных.
func RegisterOKDefined() {
	parser.Functions[okDefined] = &parser.Function{
		Name:       okDefined,
		ArgTypes:   []parser.ValueType{parser.ValueTypeMatrix},
		ReturnType: parser.ValueTypeVector,
	}
	promql.FunctionCalls[okDefined] = funcOKDefined
}

// === ok_defined(Matrix parser.ValueTypeMatrix) Vector ===
func funcOKDefined(vals []parser.Value, args parser.Expressions, enh *promql.EvalNodeHelper) promql.Vector {
	vec := vals[0].(promql.Matrix)
	for _, el := range vec {
		var value float64 = 1
		if isNone(takeLast(el), enh) {
			value = 0
		}
		enh.Out = append(enh.Out, promql.Sample{
			Metric: enh.DropMetricName(el.Metric),
			Point:  promql.Point{V: value},
		})
	}
	return enh.Out
}
