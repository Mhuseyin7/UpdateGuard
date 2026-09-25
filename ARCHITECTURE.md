# Architecture

## Trust boundaries

The control plane owns policy, plans, job transitions, audit events, and persistence. It does not receive the Docker socket. The agent owns a Docker Engine connection and accepts only authenticated, service-scoped commands. It never accepts arbitrary Docker API paths or raw compose/shell commands.

For the single-node deployment, the agent runs on the Docker host and is internal-only. Remote hosts use enrollment tokens only once, then mTLS client identities with rotation. Docker TCP must remain disabled/unpublished.

## Persistence model

`Store` is the persistence seam. The intended adapters are SQLite in WAL mode for a single owner and PostgreSQL for multi-user/multi-host use. Update jobs, snapshots, locks, policy changes, and audit records must be committed transactionally. Registry credential ciphertext belongs in a separate table with an application-managed key, never in a job payload.

## Update transaction

1. Acquire a deterministic `host/stack/service` lock and reject overlapping stack work.
2. Resolve and save both local and remote immutable digests; hash compose configuration.
3. Preflight disk, image platform, volumes, dependency health, and registry reachability.
4. Execute only a configured restricted backup integration.
5. Pull the target digest, recreate only the requested Compose service, and retain the old image.
6. Check Docker health, TCP/HTTP/custom checks, then observe restarts/exits/log patterns for the stabilization window.
7. On failure, recreate from the saved digest/configuration and verify recovery. Report `ROLLBACK_FAILED` without ambiguity if recovery fails.

## API

The public API namespace is `/api/v1` and will expose agents, hosts, stacks, services, updates, policies, registries, notifications, and audit. Mutating routes require authenticated RBAC authorization and an audit event. The control-plane probes are `/health` and `/ready`.
