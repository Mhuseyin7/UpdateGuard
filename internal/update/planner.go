package update

import (
	"fmt"
	"github.com/updateguard/updateguard/internal/domain"
	"strings"
)

type Plan struct {
	ServiceID, CurrentDigest, TargetDigest string
	Steps                                  []string
	RequiresApproval                       bool
	Warnings                               []string
}

// BuildPlan refuses an unsafe mutable-tag-only rollback path.
func BuildPlan(s domain.Service, current, target domain.Snapshot, verify domain.Verification) (Plan, error) {
	if current.Digest == "" {
		return Plan{}, fmt.Errorf("current image digest is required for rollback")
	}
	if target.Digest == "" {
		return Plan{}, fmt.Errorf("target image digest is required")
	}
	p := Plan{ServiceID: s.ID, CurrentDigest: current.Digest, TargetDigest: target.Digest, RequiresApproval: s.Policy == domain.Manual || s.Policy == domain.NotifyOnly}
	p.Steps = []string{"preflight: target image, disk, volumes and dependencies", "record immutable state snapshot", "run configured restricted backup hook", "pull target by digest", "recreate only this compose service", "verify Docker/HTTP health", fmt.Sprintf("observe stabilization for %s", verify.Stabilization), "retain previous known-good image"}
	if strings.HasSuffix(strings.Split(s.Image, ":")[len(strings.Split(s.Image, ":"))-1], "latest") {
		p.Warnings = append(p.Warnings, "floating tag: target is pinned to its resolved digest before deployment")
	}
	return p, nil
}
