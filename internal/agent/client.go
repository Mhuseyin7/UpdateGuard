// Package agent is the narrow control-plane boundary; it is deliberately not a Docker client.
package agent

import (
	"context"
	"github.com/updateguard/updateguard/internal/domain"
)

type Client interface {
	Discover(context.Context) ([]domain.Service, error)
	Snapshot(context.Context, string) (domain.Snapshot, error)
	Preflight(context.Context, string, domain.Snapshot) error
	Pull(context.Context, string) error
	Recreate(context.Context, string, domain.Snapshot) error
	Verify(context.Context, string, domain.Verification) error
	Rollback(context.Context, string, domain.Snapshot) error
	Logs(context.Context, string, int) ([]string, error)
}
