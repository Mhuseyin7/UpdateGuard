# Contributing

Use Go formatting, keep the state-machine transition table exhaustive, and add a test for every safety-sensitive behavior. Do not add a Docker operation to the agent without an authorization test and a threat-model update.

Before opening a pull request:

```bash
go test ./...
cd web && npm run typecheck && npm run build
```

Changes to deployment safety need an integration fixture covering both healthy deployment and failed verification with a successful rollback.
