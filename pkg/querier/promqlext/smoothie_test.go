package promqlext_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/grafana/mimir/pkg/querier/promqlext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/timestamp"
	"github.com/prometheus/prometheus/model/value"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/util/teststorage"
)

func TestSmoothie(t *testing.T) {
	storage := teststorage.New(t)
	a := storage.Appender(context.Background())
	ls := labels.FromStrings(labels.MetricName, "should_be_kept", "key", "value")
	minute := int64(time.Minute / time.Millisecond)
	_, _ = a.Append(0, ls, 0*minute, 0)
	_, _ = a.Append(0, ls, 1*minute, 1)
	_, _ = a.Append(0, ls, 2*minute, 2)
	_, _ = a.Append(0, ls, 3*minute, 3)
	_, _ = a.Append(0, ls, 4*minute, 4)
	_, _ = a.Append(0, ls, 5*minute, math.Float64frombits(value.StaleNaN))
	_ = a.Commit()

	promqlext.RegisterOKSmoothie()
	opts := promql.EngineOpts{
		Logger:        nil,
		Reg:           nil,
		MaxSamples:    10000,
		Timeout:       10 * time.Second,
		LookbackDelta: 5 * time.Minute,
	}
	engine := promql.NewEngine(opts)

	query, err := engine.NewRangeQuery(
		storage,
		nil,
		"ok_smoothie(should_be_kept[5m])",
		timestamp.Time(0),
		timestamp.Time(10*minute),
		time.Duration(time.Minute),
	)
	require.NoError(t, err)

	result := query.Exec(context.Background())
	require.NoError(t, result.Err)
	m := result.Value.(promql.Matrix)
	require.Len(t, m, 1)
	assert.Equal(t, ls, m[0].Metric)
	assert.Len(t, m[0].Points, 5)
	assert.EqualValues(t, m[0].Points[0].T, 0)
	assert.EqualValues(t, m[0].Points[0].V, 0)
	assert.EqualValues(t, m[0].Points[1].T, 1*minute)
	assert.EqualValues(t, m[0].Points[1].V, 0.5)
	assert.EqualValues(t, m[0].Points[2].T, 2*minute)
	assert.EqualValues(t, m[0].Points[2].V, 1)
	assert.EqualValues(t, m[0].Points[3].T, 3*minute)
	assert.EqualValues(t, m[0].Points[3].V, 1.5)
	assert.EqualValues(t, m[0].Points[4].T, 4*minute)
	assert.EqualValues(t, m[0].Points[4].V, 2)
}
