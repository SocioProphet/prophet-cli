# prophet command map

## Bootstrap façade namespace
- `prophet bootstrap doctor` -> `sourceos-bootstrap doctor`
- `prophet bootstrap login` -> `sourceos-bootstrap login`
- `prophet bootstrap build` -> `sourceos-bootstrap build`
- `prophet bootstrap validate` -> `sourceos-bootstrap validate`
- `prophet bootstrap verify` -> `sourceos-bootstrap verify`

## Agent overlay namespace
- `prophet ask ...`
- `prophet plan ...`
- `prophet agent run ...`
- `prophet mcp serve`
- `prophet mcp doctor`

## Git workflow: merge current branch into `main` without dropping content
1. `git checkout main`
2. `git merge --no-ff work`
3. Resolve conflicts by keeping both sides where needed, then run `git add <file>`
4. `git commit` (only needed if conflict resolution occurred)

Use `--no-ff` so history keeps an explicit merge commit, which helps preserve branch context.
