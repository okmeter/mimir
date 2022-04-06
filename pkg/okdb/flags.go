package okdb

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

// FlagValue variants of flag values
type FlagValue int8

const (
	// Never use feature
	Never FlagValue = iota
	// Try use feature
	Try
	// Always use feature
	Always
)

// ParseFlag parse string as FlagValue.String() result
func ParseFlag(s string) FlagValue {
	switch strings.ToLower(s) {
	case "try":
		return Try
	case "always":
		return Always
	default:
		return Never
	}
}

// String implements fmt.Stringer
func (value FlagValue) String() string {
	switch value {
	case Try:
		return "try"
	case Always:
		return "always"
	default:
		return "never"
	}
}

type useOkDBFlag struct{}

// UseOkDBFromContext get flag UseOkDB from context
func UseOkDBFromContext(ctx context.Context) FlagValue {
	value, ok := ctx.Value(useOkDBFlag{}).(FlagValue)
	if !ok {
		return Never
	}
	return value
}

const (
	// UseOkDBHTTPHeader http header with flag
	UseOkDBHTTPHeader = "X-Use-OkDB"
	// UseOkDBGRPCMetaKey grpc metadata key with flag
	UseOkDBGRPCMetaKey = "use_ok_db"
)

// ExtractUseOkDBFromHTTPHeaderMW add flag UseOkDB to context from header
func ExtractUseOkDBFromHTTPHeaderMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, useOkDBFlag{}, ParseFlag(r.Header.Get(UseOkDBHTTPHeader)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type fnRoundTripper func(*http.Request) (*http.Response, error)

func (fn fnRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

// ExtractUseOkDBFromHTTPHeaderTW add flag UseOkDB to context from header
func ExtractUseOkDBFromHTTPHeaderTW(next http.RoundTripper) http.RoundTripper {
	return fnRoundTripper(func(r *http.Request) (*http.Response, error) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, useOkDBFlag{}, ParseFlag(r.Header.Get(UseOkDBHTTPHeader)))
		return next.RoundTrip(r.WithContext(ctx))
	})
}

// InjectUseOkDBIntoHTTPRequest add flag UseOkDB to request headers from context
func InjectUseOkDBIntoHTTPRequest(ctx context.Context, r *http.Request) {
	r.Header.Set(UseOkDBHTTPHeader, UseOkDBFromContext(ctx).String())
}

// InjectUseOkDBIntoGRPCMeta add flag UseOkDB to grpc outgoing metadata from context
func InjectUseOkDBIntoGRPCMeta(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, UseOkDBGRPCMetaKey, UseOkDBFromContext(ctx).String())
}

// ExtractUseOkDBFromGRPCMeta get flag UseOkDB from grpc incomming metadata
func ExtractUseOkDBFromGRPCMeta(ctx context.Context) FlagValue {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return Never
	}

	values := meta.Get(UseOkDBGRPCMetaKey)
	if len(values) != 1 {
		return Never
	}

	return ParseFlag(values[0])
}
