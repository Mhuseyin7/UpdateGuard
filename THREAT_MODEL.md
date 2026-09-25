# Threat model

| Asset | Primary risk | Required control |
|---|---|---|
| Docker daemon | remote takeover | agent-only socket access; operation allowlist; no public agent port |
| Registry credentials | disclosure | envelope encryption, masking, no API read-back |
| Workloads and volumes | destructive update | immutable snapshot, health gates, rollback; no automatic volume delete/prune |
| Update history | tampering | append-only audit trail, authenticated RBAC, durable database |
| Remote agent | impersonation | one-time enrollment plus rotated mTLS identity |
| Backup hook | command injection | declarative trusted integrations; no arbitrary shell input |

Out of scope: a compromised Docker daemon already has host-level authority; UpdateGuard must make that trust explicit rather than hiding it behind a dashboard.
