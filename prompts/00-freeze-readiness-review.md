MODEL: claude-opus-class or fable   EFFORT: xhigh   TOOL: Codex CLI (`codex resume <thread>`) or Claude Code
WHY: adversarial, read-only review of the contract set before freeze. Repeat until a round finds nothing.

You are reviewing a FROZEN or about-to-freeze contract set, not a draft feature list.
Do not invent product requirements. Do not reopen Held items unless the text newly
contradicts them. Read-only: change no files.

Read ALL of, fully: README.md · docs/quality-bar.md · docs/contract.md ·
docs/contract-grants.md · docs/contract-boundaries.md · docs/deltas.md ·
docs/decisions.md · docs/threat-model.md · docs/onboarding.md · docs/future.md ·
docs/open-questions.md · seed/PROVENANCE.md · tools/contract-lint.py (run it: it must
pass before you start; if it fails, that is finding #1).
Cross-read the inherited contract: ~/project/mcp-sso/docs/contracts/ (§05–§19) at the
SHA named in docs/contract.md, and mcp-sso/fixtures/.

TAG SEMANTICS: `[O:date]` decided — verify each has a row in docs/decisions.md;
`[R]` confirmed default; `[D ← rule]` derived — verify the named rule exists and is
`[O]`; `[S:mcp-sso §N]` inherited — verify the clause says what is claimed; `[#]`
resolves in B8. Anything untagged that reads normative is a Pretend-decided finding.

For every never and every operation-table row: can a hostile test be written that a
WRONG build fails? Walk these classes against docs AND fixtures: the title lies; fake
id instead of a real one from the other room; the test catches its own failure; any
error counts as the check; every test is the same person; sentence and table
disagree; "who am I" demands a profile; the test login is thinner than the contract;
garbage is stored as real; a number or never appears nobody chose; one action of a
family tested, siblings not.

Output, ranked within each section, mechanical vs design separated:
Holes / Collisions / Pretend-decided / Suggestions / Held.
State explicitly which nevers or rows still lack a testable refusal. Cite file:line.
Do not invent remaining opens. Do not write code.
