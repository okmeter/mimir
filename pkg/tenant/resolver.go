// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/tenant/resolver.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package tenant

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/tenant"
	"github.com/weaveworks/common/user"
)

type MultiStarResolver struct {
	*tenant.MultiResolver
	adminTenant string
	logger      log.Logger
	cancel      func()
	allTenants  *atomic.Value // []string
}

func NewMultiStarResolver(
	ctx context.Context,
	adminTenant string,
	interval time.Duration,
	logger log.Logger,
	fetchAllTenants func(context.Context) ([]string, error),
) *MultiStarResolver {
	t := &MultiStarResolver{
		MultiResolver: tenant.NewMultiResolver(),
		adminTenant:   adminTenant,
		logger:        logger,
		allTenants:    new(atomic.Value),
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel
	t.allTenants.Store([]string{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			if tenants, err := fetchAllTenants(ctx); err == nil {
				level.Info(logger).Log("msg", "tenant resolver successful load tenants", "tenants", len(tenants))
				t.allTenants.Store(tenants)
			} else {
				level.Error(logger).Log("msg", "tenant resolver failed load tenants", "error", err)
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
	return t
}

func (t *MultiStarResolver) TenantIDs(ctx context.Context) ([]string, error) {
	orgID, err := user.ExtractOrgID(ctx)
	if err != nil {
		return nil, err
	}
	if orgID == t.adminTenant {
		tenants := t.allTenants.Load().([]string)
		level.Debug(t.logger).Log("msg", "god mode enabled", "tenants", len(tenants))
		return tenants, nil
	}
	return t.MultiResolver.TenantIDs(ctx)
}
