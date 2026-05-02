# prophet-cli

Façade repo for Prophet command surface and SourceOS bootstrap delegation.

`prophet` is the stable operator-facing command. It does not duplicate implementation logic that belongs in SourceOS or AgentTerm repositories.

## Delegation model

| Prophet command | Delegate | Owning repo |
|---|---|---|
| `prophet sourceos agent-machine ...` | `sourceosctl agent-machine ...` | `SourceOS-Linux/sourceos-devtools` |
| `prophet sourceos office ...` | `sourceosctl office ...` | `SourceOS-Linux/sourceos-devtools` |
| `prophet sourceos agent-term ...` | `agent-term ...` | `SourceOS-Linux/agent-term` |

## Examples

```bash
prophet sourceos agent-machine mounts plan
prophet sourceos agent-machine mounts init --dry-run
prophet sourceos office doctor
prophet sourceos office plan --artifact-type slide-deck --format pptx --title "Demo Deck"
prophet sourceos office generate --dry-run --artifact-type document --format docx --title "Demo Report"
prophet sourceos office convert ./example.docx --to pdf --dry-run
prophet sourceos agent-term office create-deck '!prophet-workspace' --workroom workroom-demo-0001 --title 'Demo Briefing Deck'
```

## Install path

The intended Mac-first install path is:

```bash
brew tap SocioProphet/prophet
brew install prophet-cli
brew install sourceos-devtools
brew install agent-term
```

`prophet-cli` only provides the `prophet` facade. The delegated binaries must also be installed or available on `PATH`.

## Boundary

This repo does not own:

- Agent Machine implementation;
- Office generation/conversion engines;
- LibreOffice, Collabora, ONLYOFFICE, Microsoft Graph, or Google Workspace adapters;
- AgentTerm event log semantics;
- AgentPlane evidence contracts;
- Agent Registry grants;
- Homebrew formulae.

Those remain in their owning repositories.
