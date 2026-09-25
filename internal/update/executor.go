package update

import (
 "context"
 "fmt"
 "time"
 "github.com/updateguard/updateguard/internal/agent"
 "github.com/updateguard/updateguard/internal/domain"
 "github.com/updateguard/updateguard/internal/store"
)

// Executor serializes a job's safe sequence. Any failure after deployment tries its digest-pinned rollback.
type Executor struct { Agent agent.Client; Store store.Store; Now func() time.Time }
func (e Executor) move(ctx context.Context, j *domain.UpdateJob, to domain.UpdateState) error {
 if err:=Transition(j,to);err!=nil{return err}; j.UpdatedAt=e.Now(); return e.Store.SaveJob(ctx,*j)
}
func (e Executor) audit(ctx context.Context, action, detail string) { _=e.Store.AppendAudit(ctx,domain.AuditEvent{At:e.Now(),Actor:"system",Action:action,Resource:"update",Detail:detail}) }
func (e Executor) Run(ctx context.Context, j domain.UpdateJob) domain.UpdateJob {
 fail := func(err error) domain.UpdateJob { j.Failure=err.Error(); _=e.move(ctx,&j,domain.RollingBack); e.audit(ctx,"ROLLBACK_STARTED",j.ID); if rb:=e.Agent.Rollback(ctx,j.ServiceID,j.Current); rb!=nil { j.Failure += "; rollback: "+rb.Error(); _=e.move(ctx,&j,domain.RollbackFailed); e.audit(ctx,"ROLLBACK_FAILED",j.ID); return j }; _=e.move(ctx,&j,domain.RolledBack);e.audit(ctx,"ROLLBACK_SUCCESS",j.ID);return j }
 if err:=e.move(ctx,&j,domain.Precheck);err!=nil{return j}
 if err:=e.Agent.Preflight(ctx,j.ServiceID,j.Target);err!=nil { j.Failure=err.Error(); _=e.move(ctx,&j,domain.Failed); e.audit(ctx,"UPDATE_FAILED",j.ID);return j }
 if err:=e.move(ctx,&j,domain.BackingUp);err!=nil{return j}
 // Backup hooks run exclusively in the agent's configured restricted runner; no arbitrary command is accepted here.
 if err:=e.move(ctx,&j,domain.Pulling);err!=nil{return j}
 // A pull failure has not modified the running workload, so it is a safe failure—not a rollback.
 if err:=e.Agent.Pull(ctx,j.Target.Digest);err!=nil { j.Failure="pull: "+err.Error(); _=e.move(ctx,&j,domain.Failed); e.audit(ctx,"UPDATE_FAILED",j.ID); return j }
 if err:=e.move(ctx,&j,domain.Deploying);err!=nil{return fail(err)}
 if err:=e.Agent.Recreate(ctx,j.ServiceID,j.Target);err!=nil{return fail(fmt.Errorf("deploy: %w",err))}
 if err:=e.move(ctx,&j,domain.Verifying);err!=nil{return fail(err)}
 if err:=e.Agent.Verify(ctx,j.ServiceID,j.Verification);err!=nil{return fail(fmt.Errorf("verify: %w",err))}
 if err:=e.move(ctx,&j,domain.Stabilizing);err!=nil{return fail(err)}
 // The agent's Verify endpoint observes restarts, exits and health for the complete configured window.
 if err:=e.Agent.Verify(ctx,j.ServiceID,j.Verification);err!=nil{return fail(fmt.Errorf("stabilization: %w",err))}
 _=e.move(ctx,&j,domain.Success);e.audit(ctx,"UPDATE_SUCCESS",j.ID);return j
}
