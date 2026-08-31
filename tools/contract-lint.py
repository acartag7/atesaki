#!/usr/bin/env python3
"""Contract coverage lint (open question #48).

Checks the contract set for the defect classes review rounds kept discovering:
dangling section references, audit reasons without producers, state edges without an
operation, owner tags without a ledger receipt, optional fields nothing sets, stale
section names. Reviews should verify, not discover — this makes that true for the docs.

Read-only; exits 1 on any finding. Pure standard library.
"""
from __future__ import annotations
import re, sys, pathlib

ROOT = pathlib.Path(__file__).resolve().parents[1]
DOCS = ROOT / "docs"
files = {p.name: p.read_text() for p in sorted(DOCS.glob("*.md"))}
files["seed/PROVENANCE.md"] = (ROOT / "seed" / "PROVENANCE.md").read_text()
files["README.md"] = (ROOT / "README.md").read_text()
findings: list[str] = []
def bad(msg: str) -> None: findings.append(msg)

grants = files["contract-grants.md"]; bounds = files["contract-boundaries.md"]
deltas = files["deltas.md"]; oq = files["open-questions.md"]; ledger = files["decisions.md"]

# ---- headings / ids that exist
g_heads = set(re.findall(r"^## G(\d+)\.", grants, re.M))
b_heads = set(re.findall(r"^## B(\d+)\.", bounds, re.M))
d_rows = set(re.findall(r"^\| (D\d+[abc]?) ", deltas, re.M))
op_rows = [l for l in grants.splitlines() if re.match(r"^\| (A\d+[′″ab]?|E\d) \|", l)]
op_ids = {re.match(r"^\| (A\d+[′″ab]?|E\d) \|", l).group(1) for l in op_rows}
oq_ids = set(re.findall(r"^(\d+)\.", oq, re.M)) | set(re.findall(r"^    ([a-i])\.", oq, re.M))

# ---- A. dangling references
for name, text in files.items():
    if name == "decisions.md": continue
    for m in re.finditer(r"(?<![A-Za-z/])G(\d+)\b", text):
        if m.group(1) not in g_heads: bad(f"{name}: reference G{m.group(1)} has no heading in contract-grants.md")
    for m in re.finditer(r"(?<![A-Za-z/])B(\d+)\b", text):
        if m.group(1) not in b_heads: bad(f"{name}: reference B{m.group(1)} has no heading in contract-boundaries.md")
    for m in re.finditer(r"(?<![A-Za-z/])(D\d+[abc]?)\b", text):
        if name == "deltas.md" and text[m.start()-2:m.start()] == "| ": continue
        if m.group(1) not in d_rows: bad(f"{name}: reference {m.group(1)} has no row in deltas.md")
    for m in re.finditer(r"(?<![A-Za-z/])(A\d+[′″ab]?)\b", text):
        if name != "contract-grants.md" or text[m.start()-2:m.start()] != "| ":
            if m.group(1) not in op_ids: bad(f"{name}: reference {m.group(1)} is not an operation row")
    for m in re.finditer(r"#(\d{1,2})\b", text):
        if m.group(1) not in oq_ids: bad(f"{name}: open-question reference #{m.group(1)} does not exist")
    for stale in ("§3.2", "§3.3", "GrantV0 pending", "G9)" if False else "\x00"):
        if stale in text and name != "open-questions.md": bad(f"{name}: stale reference {stale!r}")

# ---- B. audit reasons: B7 set vs G6 producers
def reasons(block: str) -> set[str]: return set(re.findall(r"`([a-z_]+)`", block))
m_d = re.search(r"\*\*D:\*\*(.*?)\n\*\*F:\*\*", bounds, re.S); m_f = re.search(r"\*\*F:\*\*(.*?)\n\n", bounds, re.S)
if not (m_d and m_f): bad("contract-boundaries.md: B7 D:/F: reason lists not found"); D_set, F_set = set(), set()
else: D_set, F_set = reasons(m_d.group(1)), reasons(m_f.group(1))
if m_d and m_f:
    for blk,name in ((m_d.group(1),"D"),(m_f.group(1),"F")):
        toks=re.findall(r"`([a-z_]+)`", blk)
        for t in set(toks):
            if toks.count(t)>1: bad(f"contract-boundaries.md: reason `{t}` listed {toks.count(t)}x in B7 {name} set")
    both=D_set & F_set
    for t in both: bad(f"contract-boundaries.md: reason `{t}` in BOTH B7 classes")
used: dict[str, set[str]] = {}
for l in op_rows:
    cells = [c.strip() for c in l.strip().strip("|").split("|")]
    if len(cells) < 6: bad(f"contract-grants.md: operation row malformed: {l[:40]}"); continue
    ev = cells[5]
    before, _, after = ev.partition("flow ")
    for r in re.findall(r"`([a-z_]+)`", before): used.setdefault(r, set()).add("D")
    for r in re.findall(r"`([a-z_]+)`", after): used.setdefault(r, set()).add("F")
for r, classes in used.items():
    if r not in D_set | F_set: bad(f"contract-grants.md: event `{r}` used in G6 but absent from B7 reason set")
    elif r in D_set and "F" in classes: bad(f"reason `{r}` is class D in B7 but written as flow in G6")
    elif r in F_set and "D" in classes and r not in ("request_deduplicated",): bad(f"reason `{r}` is class F in B7 but written as durable in G6")
for r in D_set:
    if r not in used: bad(f"contract-boundaries.md: durable reason `{r}` has no producing operation row")

# ---- C. G5 state edges vs G6 mutations
g5 = re.search(r"## G5\. States\n\n```\n(.*?)```", grants, re.S)
mut_text = " ".join(c.strip().strip("|").split("|")[4] for c in op_rows if len(c.split("|")) > 5)
if g5:
    states = set(re.findall(r"\b([a-z_]+)\b", g5.group(1))) - {"grant_request","preapproval","grant","interactive","machine","created","directly","in","its","outcome","state"}
    for st in sorted(states):
        if f"`{st}`" not in mut_text and st not in ("resolved_",): bad(f"G5 state `{st}` is never produced by a G6 mutation")
else: bad("contract-grants.md: G5 state block not found")

# ---- D. owner tags vs ledger
dates = set(re.findall(r"\| (2026-\d\d-\d\d) \|", ledger))
for name, text in files.items():
    for m in re.finditer(r"\[O:(?:(2026)-)?(\d\d-\d\d)", text):
        full = f"2026-{m.group(2)}"
        if full not in dates: bad(f"{name}: [O:{full}] has no ledger row in decisions.md")

# ---- E. optional record fields nothing sets
g2 = re.search(r"## G2\. Records(.*?)## G3\.", grants, re.S)
if g2:
    g2_no_events = re.sub(r"`grant_event`.*?\n\n", "", g2.group(1), flags=re.S)
    for f in set(re.findall(r"`([a-z_]+)\?`", g2_no_events)):
        if f"`{f}`" not in mut_text and f not in grants.split("## G6.")[1].split("## G7.")[0]:
            bad(f"G2 optional field `{f}` is never set by any operation row")

print(f"{len(findings)} finding(s)")
for f in findings: print(" -", f)
sys.exit(1 if findings else 0)
