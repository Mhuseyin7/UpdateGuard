package update
import (
 "testing"
 "time"
 "github.com/updateguard/updateguard/internal/domain"
)
func TestPlanRequiresImmutableRollback(t *testing.T) {
 _, err := BuildPlan(domain.Service{ID:"api", Image:"nginx:latest"}, domain.Snapshot{}, domain.Snapshot{Digest:"sha256:new"}, domain.Verification{Stabilization:time.Minute})
 if err == nil { t.Fatal("expected missing current digest to fail") }
}
