package promqlext

import (
	"time"

	"github.com/prometheus/prometheus/promql"
)

// ActualInterval период, в течении которого считается, что серия ещё не потеряна
//
// Агенты снимают данные каждую минуту, поэтому по умолчанию минута.
const ActualInterval = int64(time.Minute / time.Millisecond)

// Выбрать последнюю точку из серии
//
// Применяется для упрощения кода обработки даннхы в оконных функциях
func takeLast(series promql.Series) promql.Point {
	return series.Points[len(series.Points)-1]
}

// Проверка точки на отсутствие в контексте текущего среза
func isNone(p promql.Point, enh *promql.EvalNodeHelper) bool {
	return enh.Ts-p.T > ActualInterval
}
