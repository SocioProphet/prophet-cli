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

### Model fabric

```bash
prophet model route --task summarize --privacy local-first
prophet guardrail test examples/policy.json examples/input.json
prophet ledger validate
prophet ledger records
prophet agent registry list
```

## Delegation split

`prophet-cli` owns facade semantics and user experience. It should delegate to product binaries where they exist:

- `sourceos-ai` for SourceOS carry commands.
- `holmes` for language intelligence commands.
- `model-router` for model/service routing commands.
- `guardrail-fabric` for guardrail policy evaluation.
- `model-governance-ledger` for model evidence, promotion, and rollback records.
- `agent-registry` for agent spec/identity/session registry commands.

Model-fabric commands may use local development repository fallbacks when a binary is not yet installed. The fallback root is `$PROPHET_DEV_ROOT` when set, otherwise `~/dev`. Missing binaries and missing local repos must be reported explicitly as `not-yet-installed`; they must not be reported as successful execution.

## Release contract

Every delegated binary should implement:

```bash
<tool> --version
<tool> doctor
<tool> self-test
<tool> emit-evidence
```
