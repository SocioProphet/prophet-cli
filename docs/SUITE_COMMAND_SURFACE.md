# Prophet Suite Command Surface

## Position

`prophet` is the façade command for SocioProphet, SourceOS, Holmes, and the functional AI product suite.

Specialized binaries may exist, but every product path must also be available through `prophet`.

## Golden path

```bash
brew tap SocioProphet/prophet
brew install prophet-cli
prophet doctor
prophet devtools doctor
prophet devtools profile apply ai-core
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

### SourceOS installer

```bash
prophet sourceos install --target m2 --channel dev
prophet sourceos enroll
prophet sourceos status
prophet sourceos rollback list
prophet sourceos rollback apply <ref>
```

### SourceOS carry

```bash
prophet sourceos carry list
prophet sourceos carry validate
prophet sourceos carry doctor
prophet sourceos carry emit-evidence
```

### SourceOS devtools

```bash
prophet devtools doctor
prophet devtools profile list
prophet devtools profile apply core,ai-core,containers
prophet devtools profile current
prophet devtools emit-evidence
```

### Functional labs

```bash
prophet lab list
prophet lab enable image video embedding nlplab
prophet lab disable video
prophet lab status
prophet lab doctor
```

Labs are selectable. They are not installed as one giant unmanaged bundle.

### Holmes

```bash
prophet holmes analyze ./document.txt
prophet holmes search "claim"
prophet holmes graph ./document.txt
prophet holmes casefile create
prophet holmes govern ./document.txt
```

### Model routing

```bash
prophet model route --task summarize --privacy local-first
prophet model route explain <route-id>
prophet model route test examples/model-route.json
```

### Guardrails

```bash
prophet guardrail test examples/policy.json examples/input.json
prophet guardrail explain <decision-id>
```

### Agents

```bash
prophet agent registry list
prophet agent registry validate examples/agent.json
prophet agent session list
prophet agent revoke <agent-id>
```

## Delegation split

`prophet-cli` owns the façade and command semantics.

Implementation is delegated:

- SourceOS install and boot: SourceOS boot/installer repos.
- SourceOS carry: `SourceOS-Linux/sourceos-model-carry`.
- SourceOS devtools: `SourceOS-Linux/sourceos-devtools`.
- Holmes: `SocioProphet/holmes`.
- Model routing: `SocioProphet/model-router`.
- Guardrails: `SocioProphet/guardrail-fabric`.
- Agent registry: `SocioProphet/agent-registry`.
- Platform runtime: `SocioProphet/prophet-platform`.

## Install contract

Every delegated binary must support:

```bash
<tool> --version
<tool> doctor
<tool> self-test
<tool> emit-evidence
```

`prophet doctor` must verify installed delegated tools and report missing optional surfaces without failing the whole suite.
