# Prophet Suite Command Surface

`prophet` is the facade command for SocioProphet, SourceOS, Holmes, and the functional AI product suite.

Specialized binaries may exist, but every product path should also be available through `prophet`.

## Golden path

```bash
brew tap SocioProphet/prophet
brew install prophet-cli
prophet doctor
prophet devtools profile list
prophet sourceos carry list
prophet holmes analyze ./document.txt
```

## Command families

### Global

```bash
prophet version
prophet doctor
prophet self-test
prophet emit-evidence
```

### SourceOS

```bash
prophet sourceos install --target m2 --channel dev
prophet sourceos carry list
prophet sourceos carry validate
prophet sourceos carry doctor
prophet sourceos carry emit-evidence
```

### Developer tools and labs

```bash
prophet devtools doctor
prophet devtools profile list
prophet devtools profile apply core,ai-core,containers
prophet lab list
prophet lab enable image video embedding nlplab
prophet lab status
```

### Holmes and language intelligence

```bash
prophet holmes analyze ./document.txt
prophet holmes search "claim"
prophet holmes graph ./document.txt
prophet holmes govern ./document.txt
```

### Model, guardrail, and agent registry

```bash
prophet model route --task summarize --privacy local-first
prophet guardrail test examples/policy.json examples/input.json
prophet agent registry list
```

## Delegation split

`prophet-cli` owns facade semantics and user experience. It should delegate to product binaries where they exist:

- `sourceos-ai` for SourceOS carry commands.
- `holmes` for language intelligence commands.
- `model-router` for model/service routing commands.
- `guardrail-fabric` for guardrail policy evaluation.
- `agent-registry` for agent spec/identity/session registry commands.

Missing delegated tools must be reported as `not-yet-wired`, not as success.

## Release contract

Every delegated binary should implement:

```bash
<tool> --version
<tool> doctor
<tool> self-test
<tool> emit-evidence
```
