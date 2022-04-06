package promqlext_test

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/mimir/pkg/querier/promqlext"

	"github.com/stretchr/testify/suite"
)

func TestZeroIfNone(t *testing.T) {
	suite.Run(t, new(ZeroIfNoneSuite))
}

type ZeroIfNoneSuite struct {
	BaseSuite
}

func (s *ZeroIfNoneSuite) SetupSuite() {
	promqlext.RegisterOKZeroIfNone()
	s.BaseSuite.SetupSuite()
}

func (s *ZeroIfNoneSuite) TestSingle() {
	s.fillStorageByStrings(
		"1 5 __name__=foo,server=foo",
	)
	expected := `{server="foo"} => 0 @[180000]`

	query, err := s.engine.NewInstantQuery(
		s.storage,
		nil,
		"ok_zero_if_none(foo[3m])",
		s.time("3m"),
	)
	s.Require().NoError(err)

	result := query.Exec(context.Background())
	s.Require().NoError(result.Err)
	s.Equal(expected, result.String())
}

func (s *ZeroIfNoneSuite) TestMinimal() {
	s.fillStorageByStrings(
		"1 5 __name__=foo,server=foo",
	)

	query, err := s.engine.NewInstantQuery(
		s.storage,
		nil,
		"ok_zero_if_none(foo[3m], 60001)",
		s.time("3m"),
	)
	s.Require().NoError(err)

	result := query.Exec(context.Background())
	s.Require().NoError(result.Err)
	s.Zero(result.String())
}

func (s *ZeroIfNoneSuite) TestMultiple() {
	s.fillStorageByStrings(
		"0 5 __name__=foo,server=apple",
		"0 5 __name__=foo,server=orange",
		// "1 5 __name__=foo,server=apple",
		"1 0 __name__=foo,server=orange",
		// "2 5 __name__=foo,server=apple",
		"2 -42 __name__=foo,server=orange",
		// "3 5 __name__=foo,server=apple",
		"3 5 __name__=foo,server=orange",
		"4 3 __name__=foo,server=apple",
	)
	query, err := s.engine.NewRangeQuery(
		s.storage,
		nil,
		"ok_zero_if_none(foo[10m])",
		s.time("2m"),
		s.time("6m"),
		time.Duration(time.Minute),
	)
	s.Require().NoError(err)

	result := query.Exec(context.Background())
	s.Require().NoError(result.Err)
	s.checkMatrix(result, map[string][]float64{
		// Minutest           2  3  4  5  6
		`{server="apple"}`:  {0, 0, 3, 3, 0},
		`{server="orange"}`: {-42, 5, 5, 0, 0},
	})
}
