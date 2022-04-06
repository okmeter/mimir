package okdb

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/middleware"
	"github.com/grafana/mimir/pkg/okdb/okdbfrontpb"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sony/gobreaker"
	"github.com/thanos-io/thanos/pkg/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// MakeClient connection pool
func MakeClient(cfg Config, logger log.Logger, reg prometheus.Registerer) (okdbfrontpb.OkStorageFrontendClient, error) {
	target := fmt.Sprintf("dns:///%s:%d", cfg.Address, cfg.GetPort())
	conn, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(cfg.GetReceiveMsgSize()),
		),
		grpc.WithChainUnaryInterceptor(
			makeCircuitBreakerInterceptor(cfg, logger, reg),
			tracing.UnaryClientInterceptor(opentracing.GlobalTracer()),
			makeMetricInterceptor(reg),
			makeTimeoutInterceptor(cfg),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", target, err)
	}
	return okdbfrontpb.NewOkStorageFrontendClient(conn), nil
}

func makeCircuitBreakerInterceptor(cfg Config, logger log.Logger, reg prometheus.Registerer) grpc.UnaryClientInterceptor {
	logger = log.With(logger, "component", "okdb-client-circuit-breaker")
	state := promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "okdb",
		Subsystem: "client",
		Name:      "circuit_breaker_state",
		Help:      "Current circuit breaker state.",
	}, []string{"state"})
	state.WithLabelValues(gobreaker.StateClosed.String()).SetToCurrentTime()
	level.Info(logger).Log("msg", "state changed", "state", gobreaker.StateClosed.String())
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "okDBFront",
		MaxRequests: 100,
		Interval:    0,
		Timeout:     cfg.GetBreakerTimeout(),
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > cfg.GetBreakerMaxFailures()
		},
		OnStateChange: func(_ string, _, to gobreaker.State) {
			state.WithLabelValues(to.String()).SetToCurrentTime()
			level.Info(logger).Log("msg", "state changed", "state", to.String())
		},
		IsSuccessful: func(err error) bool {
			return err == nil || status.Code(err) == codes.NotFound
		},
	})
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if UseOkDBFromContext(ctx) == Always {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		_, err := cb.Execute(func() (interface{}, error) {
			return nil, invoker(ctx, method, req, reply, cc, opts...)
		})
		return err
	}
}

func makeMetricInterceptor(reg prometheus.Registerer) grpc.UnaryClientInterceptor {
	duration := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "okdb",
		Subsystem: "client",
		Name:      "request_duration_seconds",
		Help:      "Time spent executing requests to the frontend.",
		Buckets:   prometheus.ExponentialBuckets(0.008, 4, 7),
	}, []string{"operation", "statusCode"})
	return middleware.PrometheusGRPCUnaryInstrumentation(duration)
}

func makeTimeoutInterceptor(cfg Config) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx, cancel := context.WithTimeout(ctx, cfg.GetRequestTimeout())
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
