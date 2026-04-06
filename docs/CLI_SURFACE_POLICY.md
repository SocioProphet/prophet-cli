# CLI surface policy

The deterministic CLI surface is canonical.

## Read verbs
- `list`
- `show`
- `describe`
- `explain`
- `find`
- `search`

## Mutation verbs
- `create`
- `update`
- `delete`
- `apply`
- `wait`
- `validate`
- `verify`

## Runtime verbs
- `status`
- `logs`
- `start`
- `stop`
- `restart`

## Bootstrap verbs
- `doctor`
- `login`
- `build`
- `fetch`
- `write`
- `info`

Agent affordances may exist, but they must resolve to deterministic tools and must not define a competing command language.
