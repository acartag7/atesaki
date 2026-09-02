MODEL: claude-opus-class or fable   EFFORT: high   TOOL: Claude Code in ~/project/atesaki-core MILESTONE: M0 (docs/roadmap.md). Serial PRs; one behavior each; merge before the next. WHY: Go code is merging on a Codex read and the implementer's local run. Nothing runs the suite, nothing protects main, the public repo has no license or disclosure channel, and Go has no package-manager release-age gate. Every later milestone assumes these.

Read first, fully: docs/roadmap.md §M0 · docs/quality-bar.md · prompts/README.md · .github/ (empty today) · go.mod · internal/config/load.go and cmd/atesaki/main.go (the two nits) · STATE.md · docs/decisions.md.

PR 1: `ci: build, vet, gofmt, race tests, govulncheck, contract lint`
- `.github/workflows/ci.yml`: on pull_request and push to main; `permissions: contents: read`; a matrix job (`ubuntu-latest` + `macos-latest`) **plus one aggregate job named `ci` that `needs` the matrix and runs with `if: always()`, failing unless every leg's result is `success`**, a job that merely `needs` a failed matrix is skipped, and GitHub counts a skipped required check as passing; branch protection requires `ci`, because matrix legs report suffixed names; Go version from go.mod; actions pinned by commit SHA (version in a comment); steps: `go build ./...`, `go vet ./...`, `test -z "$(gofmt -l .)"` (a bare `gofmt -l` exits 0), `go test -race ./...` (`-mod=readonly`), `go run golang.org/x/vuln/cmd/govulncheck@<exact version> ./...` (pinned; a fresh runner has none), `python3 tools/contract-lint.py`, `git diff --check`.
- Prove the gate bites twice: a failing test makes `ci` red (link the run), and after protection is on, a PR with a deliberately failed leg cannot be merged (screenshot or API output in the PR); then drop the commit and link the green run.
- After merge (owner runs, commands in the PR body): protection on `main` as a **ruleset**: required status check `ci`, no force-push, no deletion, linear history, require the branch to be up to date, and **bypass disabled for administrators** (GitHub lets admins bypass classic protection by default; the sole merging owner must not be able to merge red). Record the exact `gh api` call.
- Sanitization grep as a CI step from this PR on: hostnames, tenant ids, AD group names, vault paths, employer names, and references to the private evidence folder outside `docs/` evidence notes, packet 09's categories, run before every public push from now on (every PR before M8 is already public; M8 re-runs this, it does not introduce it).

PR 2: `chore: license, security policy, gitignore`
- `LICENSE`: the owner's choice (`[decide]`). mcp-sso ships MIT; Apache-2.0 adds an explicit patent grant and NOTICE handling. Present both in the PR with that tradeoff and apply the one the owner names; do not pick.
- `SECURITY.md`: supported versions (v0 pre-release: main only), how to report (GitHub private vulnerability reporting, enabled in repo settings by the owner, say so), the response window, what counts as in scope (the binary, the container, the recipe), the public-fix rule (no narration of an unreleased vulnerability in commits until the release is out).
- `.gitignore`: the built binary, `*.test`, coverage files, `.DS_Store`, editor scratch. Nothing that is already tracked.

PR 3: `chore: dependency cooldown and review`
- The floor `[decide]`: the plan assumes **15 days, majors 30** (mcp-sso's number, Atesaki is the same kind of boundary with a handful of dependencies); the house minimum is 5/30.
- Bot layer: `.github/dependabot.yml` (Dependabot, not Renovate) for `gomod` and `github-actions` with `cooldown: default-days: <floor>, semver-major-days: 30` (GitHub's built-in default has been 3 days since 2026-07-14, set ours explicitly).
- There is **no honest CI age gate for Go**: the module proxy's `.info` `Time` is the commit time of the tagged revision, not a publish time, so a freshly tagged old commit passes any age test. Do not build one. Instead: GitHub's dependency-review action on every PR that changes `go.mod`, and a PR template line that asks for the version's publish date and age as human evidence, with the advisory id when the cooldown is bypassed. Record the gap in `docs/quality-bar.md` in one sentence.
- State in the PR that `go.sum` is the integrity ledger, not a lockfile: `go.mod` plus minimal version selection pins the build list, checked with `go list -m all`.

PR 4: `fix(config): read the config file through one descriptor`
- `Load` stats the path, then reads the path again: the size cap is advisory across that gap. Fix: open once (symlinks allowed, a ConfigMap-mounted config is one), `fstat` the descriptor for the regular-file check only, read through `io.LimitReader(f, configMaxBytes+1)`, refuse if more than the cap arrived. The cap is enforced by the reader, so no size-then-read race exists to test; regression tests: a file one byte over the cap is refused, a file exactly at the cap is read, and the refusal names `B5.config-size`.

PR 5: `refactor(config): one reserved-path source`
- `cmd/atesaki/main.go` and `internal/config/validate.go` each build the reserved list. Expose it once from `config` and use it in both. A test asserts the `routes` output's `reserved` equals what collision checking used.

PR 6: `fix(config): empty port; http redirects off loopback`
- `checkHostPort` accepts a present but empty port (`gw.example.com:`, `[::1]:`), the unmirrored sibling of PR 6's `checkURL` fix; refuse it, with cases for a DNS name, IPv4, and bracketed IPv6.
- `checkOriginOrURL` (redirect allowlist entries) accepts `http://` on any host; the inherited mcp-sso §10 allows `http` only on loopback, refuse non-loopback `http` now (cases: an origin, an exact URL, loopback IPv4, loopback IPv6, a non-loopback `http`). The **full** inherited §10 entry grammar (delimiter and control checks, canonical spelling, root-slash and percent-triplet rules, length bounds) is slice 2's first PR with its `inherited` fixtures, not this one.
- The credential header-name refusal (`Host`, `Content-Length`, `Transfer-Encoding`, `Connection`, `Trailer`, `Upgrade`, `TE`, `Keep-Alive`, `Proxy-*`) is written **after** packet 14 item 5 lands the B1 sentence, contract before code, as its own small PR, one case per forbidden name.
- Every case under `testdata/refuse/` with `# expect:` naming the rule; the mirror sweep listed in the PR.

PR 7: `docs: name check`, packet 09's deliverable 1, run now: `atesaki` on GitHub (org and repo), Go module path, container registry namespace, npm and PyPI collisions (informational), the search-query family and autocomplete, a trademark screen; each free / taken / congested; a ledger row once the owner rules (open question #9). No rename is performed here.

PR 8: `docs: record slices-before-freeze; refresh STATE; reorder packets`
- `docs/decisions.md`: a 2026-09-01 row "product code proceeds slice by slice against the draft; each PR names the sections it implements" with receipt "PR 5 merged with that README line". Strike nothing; it is a new row.
- `STATE.md`: lane 1 rewritten to the serial-PR reality (PR 340 superseded), the PR 3 row corrected, the probe-server line removed, the blocked-on-owner list rebuilt from docs/roadmap.md §2.
- `prompts/README.md` already carries the roadmap order; verify it matches docs/roadmap.md §3 and fix drift.

LOCAL, NOT A PR: `git branch -d` the three merged branches; remove the dead `probe-a`/`probe-b` `[mcp_servers.*]` entries from `~/.codex/config.toml` (they error on every Codex start); leave the Claude Code local-scope probe entries until the manual authorize check in STATE.md is done.

HARD RULES: no product behavior changes beyond PR 4, PR 5, and PR 6; nothing touches docs/contract*.md; every PR verified by running (`go test -race ./...`, the workflow itself). Windows is out of v0 (B2 uses `O_NOFOLLOW` and `Stat_t`); README says "Linux and macOS" in PR 2.

DONE WHEN: CI is a required check on main; a red PR cannot be merged by the owner's own account (proven); the sanitization grep runs in CI; Dependabot cooldown and dependency review in place; the two config fixes merged with regression tests and their mirrors; the name check reported; STATE and ledger current.

REPORT: the red-then-green run links; the `gh api` protection call; the license chosen and why; anything the workflow could not run on macOS.
