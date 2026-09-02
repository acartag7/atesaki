MODEL: claude-opus-class or fable   EFFORT: high   TOOL: Claude Code in ~/project/atesaki-core
MILESTONE: M0 (docs/roadmap.md). Serial PRs; one behavior each; merge before the next.
WHY: Go code is merging on a Codex read and the implementer's local run. Nothing runs
the suite, nothing protects main, the public repo has no license or disclosure channel,
and Go has no package-manager release-age gate. Every later milestone assumes these.

Read first, fully: docs/roadmap.md §M0 · docs/quality-bar.md · prompts/README.md ·
.github/ (empty today) · go.mod · internal/config/load.go and cmd/atesaki/main.go
(the two nits) · STATE.md · docs/decisions.md.

PR 1 — `ci: build, vet, gofmt, race tests, govulncheck, contract lint`
- `.github/workflows/ci.yml`: on pull_request and push to main; `permissions:
  contents: read`; a matrix job (`ubuntu-latest` + `macos-latest`) **plus one
  aggregate job named `ci` that `needs` the matrix** — branch protection requires
  `ci`, because matrix legs report suffixed names; Go version from go.mod; actions
  pinned by commit SHA (version in a comment); steps: `go build ./...`, `go vet
  ./...`, `test -z "$(gofmt -l .)"` (a bare `gofmt -l` exits 0), `go test -race
  ./...` (`-mod=readonly`), `go run golang.org/x/vuln/cmd/govulncheck@<exact
  version> ./...` (pinned; a fresh runner has none), `python3
  tools/contract-lint.py`, `git diff --check`.
- Prove the gate bites: push one temporary commit with a failing test, link the red
  run in the PR, drop the commit, link the green run.
- After merge (owner runs, commands in the PR body): branch protection on `main` —
  required status check `ci`, no force-push, no deletion, linear history, require the
  branch to be up to date. Record the exact `gh api` call.

PR 2 — `chore: license, security policy, gitignore`
- `LICENSE`: the owner's choice. Default: the same license mcp-sso ships (read its
  LICENSE; do not assume); state the choice in the PR as `[decide]` if it differs.
- `SECURITY.md`: supported versions (v0 pre-release: main only), how to report (GitHub
  private vulnerability reporting, enabled in repo settings by the owner — say so),
  the response window, what counts as in scope (the binary, the container, the
  recipe), the public-fix rule (no narration of an unreleased vulnerability in
  commits until the release is out).
- `.gitignore`: the built binary, `*.test`, coverage files, `.DS_Store`, editor
  scratch. Nothing that is already tracked.

PR 3 — `chore: dependency cooldown in CI and the update bot`
- The floor `[decide]`: the plan assumes **15 days, majors 30** (mcp-sso's number —
  Atesaki is the same kind of boundary with a handful of dependencies); the house
  minimum is 5/30. One number, used by both layers.
- Bot layer: `.github/dependabot.yml` (Dependabot, not Renovate) for `gomod` and
  `github-actions` with `cooldown: default-days: <floor>, semver-major-days: 30`
  (GitHub's built-in default has been 3 days since 2026-07-14 — set ours
  explicitly; a bot-only floor stops nothing).
- CI layer: `tools/depage.py` (stdlib only): for every `require` in go.mod, fetch
  `https://proxy.golang.org/<escaped module>/@v/<version>.info` (module paths are
  escaped per the Go module reference: an uppercase letter becomes `!` plus its
  lowercase) and read `Time`; fail if the version is younger than the floor, or 30
  days when its major differs from the major on `origin/main` (a module new to
  go.mod has no baseline: apply 30 days); a network failure is a CI failure, never a
  pass; an exception needs a row in `tools/dependency-exceptions.json` with module,
  version, advisory id (GHSA/CVE), minimum fixing version, and date — never a global
  floor change. The proxy timestamp is a risk signal, not provenance (`go.sum` does
  not authenticate it). Test against a fake `.info` one day old (must fail) and one
  older than the floor (must pass). Wire it into `ci.yml`. State in the PR that
  `go.sum` is the integrity ledger, not a lockfile: `go.mod` plus minimal version
  selection pins the build list, checked with `go list -m all`.

PR 4 — `fix(config): read the config file through one descriptor`
- `Load` stats the path, then reads the path again: the size cap is advisory across
  that gap. Fix: open once (symlinks allowed — a ConfigMap-mounted config is one),
  `fstat` the descriptor for the regular-file check only, read through
  `io.LimitReader(f, configMaxBytes+1)`, refuse if more than the cap arrived. The cap
  is enforced by the reader, so no size-then-read race exists to test; regression
  tests: a file one byte over the cap is refused, a file exactly at the cap is read,
  and the refusal names `B5.config-size`.

PR 5 — `refactor(config): one reserved-path source`
- `cmd/atesaki/main.go` and `internal/config/validate.go` each build the reserved
  list. Expose it once from `config` and use it in both. A test asserts the `routes`
  output's `reserved` equals what collision checking used.

PR 6 — `docs: record slices-before-freeze; refresh STATE; reorder packets`
- `docs/decisions.md`: a 2026-09-01 row "product code proceeds slice by slice against
  the draft; each PR names the sections it implements" with receipt "PR 5 merged with
  that README line". Strike nothing; it is a new row.
- `STATE.md`: lane 1 rewritten to the serial-PR reality (PR 340 superseded), the PR 3
  row corrected, the probe-server line removed, the blocked-on-owner list rebuilt from
  docs/roadmap.md §2.
- `prompts/README.md` already carries the roadmap order; verify it matches
  docs/roadmap.md §3 and fix drift.

LOCAL, NOT A PR: `git branch -d` the three merged branches; remove the dead
`probe-a`/`probe-b` `[mcp_servers.*]` entries from `~/.codex/config.toml` (they error
on every Codex start); leave the Claude Code local-scope probe entries until the
manual authorize check in STATE.md is done.

HARD RULES: no product behavior changes beyond PR 4 and PR 5; nothing touches
docs/contract*.md; every PR verified by running (`go test -race ./...`, the workflow
itself). Windows is out of v0 (B2 uses `O_NOFOLLOW` and `Stat_t`); README says
"Linux and macOS" in PR 2.

DONE WHEN: CI is a required check on main; a red PR cannot merge; the age test runs in
CI; both nits merged with regression tests; STATE and ledger current.

REPORT: the red-then-green run links; the `gh api` protection call; the license
chosen and why; anything the workflow could not run on macOS.
