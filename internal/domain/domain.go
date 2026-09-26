// Package domain contains the stable types shared by the control plane and agent.
package domain

import "time"

type UpdateState string

const (
	Pending        UpdateState = "PENDING"
	Precheck       UpdateState = "PRECHECK"
	BackingUp      UpdateState = "BACKING_UP"
	Pulling        UpdateState = "PULLING"
	Deploying      UpdateState = "DEPLOYING"
	Verifying      UpdateState = "VERIFYING"
	Stabilizing    UpdateState = "STABILIZING"
	Success        UpdateState = "SUCCESS"
	Failed         UpdateState = "FAILED"
	RollingBack    UpdateState = "ROLLING_BACK"
	RolledBack     UpdateState = "ROLLED_BACK"
	RollbackFailed UpdateState = "ROLLBACK_FAILED"
	Cancelled      UpdateState = "CANCELLED"
)

type Policy string

const (
	Manual     Policy = "MANUAL"
	NotifyOnly Policy = "NOTIFY_ONLY"
	AutoPatch  Policy = "AUTO_PATCH"
	AutoMinor  Policy = "AUTO_MINOR"
	AutoAll    Policy = "AUTO_ALL"
)

type Service struct {
	ID, HostID, Stack, Name, Image, CurrentDigest, TargetDigest string
	Health                                                      string
	RestartCount                                                int
	Policy                                                      Policy
	DependsOn                                                   []string
}

type Snapshot struct {
	ImageRef, Digest, ComposeHash string
	CreatedAt                     time.Time
	Volumes, Networks             []string
	Ports                         []string
	RestartPolicy                 string
}

type Verification struct {
	DockerHealth   bool
	HTTPURL        string
	ExpectedStatus int
	BodyContains   string
	Stabilization  time.Duration
}

type UpdateJob struct {
	ID, ServiceID        string
	State                UpdateState
	Current, Target      Snapshot
	Verification         Verification
	CreatedAt, UpdatedAt time.Time
	Failure              string
}

type AuditEvent struct {
	At                              time.Time `json:"at"`
	Actor, Action, Resource, Detail string
}
