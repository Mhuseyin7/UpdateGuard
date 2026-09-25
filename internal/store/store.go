package store

import (
 "context"
 "sync"
 "github.com/updateguard/updateguard/internal/domain"
)

// Store is the persistence seam. Production adapters use SQLite (single-node) or PostgreSQL;
// the in-memory adapter makes control-plane behaviour testable without a database.
type Store interface { SaveJob(context.Context, domain.UpdateJob) error; Job(context.Context,string)(domain.UpdateJob,error); AppendAudit(context.Context,domain.AuditEvent) error; Audit(context.Context)[]domain.AuditEvent }
type Memory struct { mu sync.RWMutex; jobs map[string]domain.UpdateJob; audit []domain.AuditEvent }
func NewMemory()*Memory{return &Memory{jobs:map[string]domain.UpdateJob{}}}
func(m *Memory)SaveJob(_ context.Context,j domain.UpdateJob)error{m.mu.Lock();defer m.mu.Unlock();m.jobs[j.ID]=j;return nil}
func(m *Memory)Job(_ context.Context,id string)(domain.UpdateJob,error){m.mu.RLock();defer m.mu.RUnlock();return m.jobs[id],nil}
func(m *Memory)AppendAudit(_ context.Context,e domain.AuditEvent)error{m.mu.Lock();defer m.mu.Unlock();m.audit=append(m.audit,e);return nil}
func(m *Memory)Audit(_ context.Context)[]domain.AuditEvent{m.mu.RLock();defer m.mu.RUnlock();return append([]domain.AuditEvent(nil),m.audit...)}
