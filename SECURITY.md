# Security policy

Do not report secrets in public issues. Report vulnerabilities privately to the repository maintainers with reproduction steps, impact, and affected version. We will acknowledge reports within seven days and coordinate a fix before disclosure.

## Operational requirements

- Keep the agent on a private network; never expose its port or the Docker socket.
- Use a unique service identity per agent, mutual TLS for remote agents, and short-lived enrollment tokens.
- Encrypt registry credentials at rest; mask them in all UI, events, and logs.
- Run the UI and control plane without Docker access. The agent must expose an allowlist, not a generic proxy.
- Treat custom backup hooks as privileged integrations: fixed executable, fixed arguments, fixed working directory, and no user-controlled interpolation.
- Do not enable automated major-version updates by default.

The initial scaffold is not yet suitable for production management; see the status note in README.
