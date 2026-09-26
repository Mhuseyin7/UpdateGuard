package update

import (
	"context"
	"errors"
	"github.com/updateguard/updateguard/internal/domain"
	"github.com/updateguard/updateguard/internal/store"
	"testing"
	"time"
)

type fakeAgent struct {
	failVerify bool
	failPull   bool
	rolled     bool
}

func (f *fakeAgent) Discover(context.Context) ([]domain.Service, error) { return nil, nil }
func (f *fakeAgent) Snapshot(context.Context, string) (domain.Snapshot, error) {
	return domain.Snapshot{}, nil
}
func (f *fakeAgent) Preflight(context.Context, string, domain.Snapshot) error { return nil }
func (f *fakeAgent) Pull(context.Context, string) error {
	if f.failPull {
		return errors.New("registry unavailable")
	}
	return nil
}
func (f *fakeAgent) Recreate(context.Context, string, domain.Snapshot) error { return nil }
func (f *fakeAgent) Verify(context.Context, string, domain.Verification) error {
	if f.failVerify {
		return errors.New("unhealthy")
	}
	return nil
}
func (f *fakeAgent) Rollback(context.Context, string, domain.Snapshot) error {
	f.rolled = true
	return nil
}
func (f *fakeAgent) Logs(context.Context, string, int) ([]string, error) { return nil, nil }
func TestFailedHealthRollsBack(t *testing.T) {
	a := &fakeAgent{failVerify: true}
	s := store.NewMemory()
	e := Executor{Agent: a, Store: s, Now: time.Now}
	got := e.Run(context.Background(), domain.UpdateJob{ID: "j", ServiceID: "s", State: domain.Pending, Current: domain.Snapshot{Digest: "sha256:v1"}, Target: domain.Snapshot{Digest: "sha256:v2"}})
	if got.State != domain.RolledBack || !a.rolled {
		t.Fatalf("got %s rollback=%v", got.State, a.rolled)
	}
}
func TestPullFailureDoesNotRollback(t *testing.T) {
	a := &fakeAgent{failPull: true}
	e := Executor{Agent: a, Store: store.NewMemory(), Now: time.Now}
	got := e.Run(context.Background(), domain.UpdateJob{ID: "j", ServiceID: "s", State: domain.Pending, Current: domain.Snapshot{Digest: "sha256:v1"}, Target: domain.Snapshot{Digest: "sha256:v2"}})
	if got.State != domain.Failed || a.rolled {
		t.Fatalf("got %s rollback=%v", got.State, a.rolled)
	}
}
