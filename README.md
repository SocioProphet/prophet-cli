# prophet-cli

Façade repo for Prophet command surface and SourceOS / AgentPlane bootstrap delegation.

`prophet` is the stable operator-facing command. It does not duplicate implementation logic that belongs in SourceOS, AgentTerm, or AgentPlane repositories.

## Delegation model

| Prophet command | Delegate | Owning repo |
|---|---|---|
| `prophet sourceos local-model ...` | `sourceosctl local-model ...` | `SourceOS-Linux/sourceos-devtools` |
| `prophet sourceos network ...` | `sourceosctl network ...` | `SourceOS-Linux/sourceos-devtools` |
| `prophet sourceos native-assistant ...` | `sourceosctl native-assistant ...` | `SourceOS-Linux/sourceos-devtools` |
| `prophet sourceos agent-machine ...` | `sourceosctl agent-machine ...` | `SourceOS-Linux/sourceos-devtools` |
| `prophet sourceos office ...` | `sourceosctl office ...` | `SourceOS-Linux/sourceos-devtools` |
| `prophet sourceos agent-term ...` | `agent-term ...` | `SourceOS-Linux/agent-term` |
| `prophet agentplane ...` | `sp-run ...` | `SocioProphet/agentplane` |
| `prophet governed-runner ...` | `sp-run ...` | `SocioProphet/agentplane` |

## Examples

```bash
prophet sourceos local-model doctor
prophet sourceos local-model profiles
prophet sourceos local-model plan --profile local-llama32-1b
prophet sourceos local-model route --task-class office-assist
prophet sourceos network doctor
prophet sourceos network plan --enterprise --mesh --allow-listed --destination models.enterprise.example
prophet sourceos network provider --provider-class openai-compatible --owner user
prophet sourceos native-assistant plan --operation open-workroom
prophet sourceos native-assistant plan --operation create-office-artifact
prophet sourceos agent-machine mounts plan
prophet sourceos agent-machine mounts init --dry-run
prophet sourceos office doctor
prophet sourceos office plan --artifact-type slide-deck --format pptx --title "Demo Deck"
prophet sourceos office generate --dry-run --artifact-type document --format docx --title "Demo Report"
prophet sourceos office convert ./example.docx --to pdf --dry-run
prophet sourceos agent-term office create-deck '!prophet-workspace' --workroom workroom-demo-0001 --title 'Demo Briefing Deck'
prophet agentplane doctor
prophet agentplane preflight ./governed-run-contract.json
prophet agentplane admit ./governed-run-contract.json --preflight ./preflight-receipt.json --authority-state ./agent-authority-current-state.json --projected-cost-usd 0.25
prophet agentplane dossier ./.socioprophet/runs/governed-run-alpha-001
prophet agentplane validate-dossier ./run-dossier.json
prophet governed-runner doctor
prophet governed-runner smoke --output-dir ./.socioprophet/smoke/governed-runner
prophet governed-runner list --runs-root ./.socioprophet/smoke/governed-runner
prophet governed-runner status ./.socioprophet/smoke/governed-runner/run
prophet governed-runner inspect ./.socioprophet/smoke/governed-runner/run
prophet governed-runner tool list-tools
prophet governed-runner tool call governed_runner.doctor --args-json '{}'
prophet governed-runner tool call governed_runner.smoke --args-json '{"output_dir":".socioprophet/smoke/governed-runner"}'
prophet governed-runner preflight ./governed-run-contract.json
prophet governed-runner admit ./governed-run-contract.json --preflight ./preflight-receipt.json --authority-state ./agent-authority-current-state.json --projected-cost-usd 0.25
prophet governed-runner dossier ./.socioprophet/runs/governed-run-alpha-001
prophet governed-runner validate-dossier ./run-dossier.json
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

For AgentPlane governed-runner commands, install or expose the AgentPlane-owned `sp-run` delegate from `SocioProphet/agentplane`:

```bash
python3 -m pip install -e /path/to/agentplane
sp-run doctor
sp-run smoke --output-dir ./.socioprophet/smoke/governed-runner
sp-run list --runs-root ./.socioprophet/smoke/governed-runner
sp-run tool list-tools
prophet agentplane doctor
prophet governed-runner doctor
prophet governed-runner smoke --output-dir ./.socioprophet/smoke/governed-runner
prophet governed-runner list --runs-root ./.socioprophet/smoke/governed-runner
prophet governed-runner tool list-tools
```

## Boundary

This repo does not own:

- Local Model Door runtime logic;
- Network Door / Firewall Door / Mesh Door implementation;
- BYOM provider connectivity or credentials;
- native assistant runtime adapters, Siri/App Intents/Shortcuts integration, Android intents, Windows shell APIs, browser extension transports, or MCP/native transports;
- model weights, model downloads, or inference;
- personal tuning or personalization governance;
- Agent Machine implementation;
- Office generation/conversion engines;
- LibreOffice, Collabora, ONLYOFFICE, Microsoft Graph, or Google Workspace adapters;
- AgentTerm event log semantics;
- AgentPlane evidence contracts;
- AgentPlane governed-runner implementation;
- Agent Registry grants or authority state;
- Homebrew formulae.

Those remain in their owning repositories.

`prophet agentplane ...` and `prophet governed-runner ...` are facades over `sp-run ...`; the implementation remains in `SocioProphet/agentplane`.
