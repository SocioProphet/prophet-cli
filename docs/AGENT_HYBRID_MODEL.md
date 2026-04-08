# Agent hybrid model

The repository follows a hybrid model with three surfaces:

1. deterministic CLI
2. agent assist
3. agent execute

## Deterministic CLI
The deterministic CLI is the canonical, documented surface.

## Agent assist
Read-mostly planning and explanation over deterministic tools.

## Agent execute
Approval-gated execution that still resolves to deterministic tools.

## Guardrails
- no bootstrap business logic in this repo
- no silent plugin installation
- no automatic project MCP loading in untrusted workspaces
- no competing command language that bypasses deterministic verbs
