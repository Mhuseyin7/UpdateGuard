// Package update enforces the safe update lifecycle. No caller may bypass transitions.
package update

import (
	"fmt"
	"github.com/updateguard/updateguard/internal/domain"
)

var allowed = map[domain.UpdateState]map[domain.UpdateState]bool{
	domain.Pending:     {domain.Precheck: true, domain.Cancelled: true},
	domain.Precheck:    {domain.BackingUp: true, domain.Failed: true, domain.Cancelled: true},
	domain.BackingUp:   {domain.Pulling: true, domain.Failed: true},
	domain.Pulling:     {domain.Deploying: true, domain.Failed: true},
	domain.Deploying:   {domain.Verifying: true, domain.RollingBack: true, domain.Failed: true},
	domain.Verifying:   {domain.Stabilizing: true, domain.RollingBack: true, domain.Failed: true},
	domain.Stabilizing: {domain.Success: true, domain.RollingBack: true},
	domain.Failed:      {domain.RollingBack: true},
	domain.RollingBack: {domain.RolledBack: true, domain.RollbackFailed: true},
}

func Transition(job *domain.UpdateJob, next domain.UpdateState) error {
	if !allowed[job.State][next] {
		return fmt.Errorf("invalid update transition %s -> %s", job.State, next)
	}
	job.State = next
	return nil
}

func IsTerminal(s domain.UpdateState) bool {
	return s == domain.Success || s == domain.RolledBack || s == domain.RollbackFailed || s == domain.Cancelled
}
