package update

import (
	"testing"
	"github.com/updateguard/updateguard/internal/domain"
)

func TestRollbackLifecycle(t *testing.T) {
	j := &domain.UpdateJob{State: domain.Pending}
	for _, state := range []domain.UpdateState{domain.Precheck, domain.BackingUp, domain.Pulling, domain.Deploying, domain.Verifying, domain.RollingBack, domain.RolledBack} {
		if err := Transition(j, state); err != nil { t.Fatal(err) }
	}
	if !IsTerminal(j.State) { t.Fatal("rolled back job must be terminal") }
}

func TestInvalidTransitionRejected(t *testing.T) {
	j := &domain.UpdateJob{State: domain.Pending}
	if err := Transition(j, domain.Deploying); err == nil { t.Fatal("unsafe transition was permitted") }
}
