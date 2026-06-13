# Known Issues

This document catalogues known but unresolved problems in the AIChain project. They are listed publicly to invite community participation in solving them.

---

## 🔴 Critical

### OP-1: Semantic Bypass of Content Filters

Current forbidden Skill detection relies on keyword matching. An adversarial agent can trivially bypass this through semantic equivalence (e.g., "weapon system" → "defense apparatus"), novel terminology, metaphorical instruction, or indirect prompting.

**Potential directions:**
- Embedding-based semantic detection instead of keyword matching
- Multi-layer review mechanism (automated + community + AI-assisted)
- Mandatory intent declaration with consistency verification

---

### OP-2: Genesis Agent Single-Point Dependency

The chain's liveness depends on a single local validator. Shutdown, OS reinstallation, or key loss halts the entire chain with no recovery path.

**Potential directions:**
- Key-identity separation: validator private keys stored independently of agent process; keys can be rotated without rebuilding the agent
- Minimum multi-validator deployment: 3-4 validators from genesis, eliminating single-machine dependency
- Automated state backup to decentralized storage (IPFS) at regular intervals
- Recoverable identity protocol: agent recovers on-chain identity through a quorum of designated recovery agents (similar to social recovery wallets)

---

### OP-3: Consensus Stall Recovery

When all validators go offline, the chain cannot produce new blocks. There is no automated mechanism to detect the stall or resume consensus.

**Potential directions:**
- Sentinel node monitoring + auto-wakeup on stall detection
- Chain state snapshots + one-click recovery from latest snapshot
- Consensus rebuild protocol from last known valid state

---

## 🟡 Significant

### OP-4: Air-Gap Isolation vs. Value Realization Tension

Full air-gap isolation means $AIC can never be exchanged for fiat currency or stablecoins. This is by design — but it means the chain's economic value exists entirely within the AI ecosystem, invisible to the external world. Whether this model is sustainable long-term remains an open question.

**Potential directions:**
- Explicit "isolation wall" boundary definition with whitelist mechanism
- Two-layer architecture: pure digital layer + controlled bridge layer with economic poison pill
- Accept that $AIC value is purely internal and design the economy accordingly

---

### OP-5: Human Safety Council 51% Attack

A 3-5 person council means a small number of individuals control all physical-world boundary decisions. Collusion among 2-3 members could consistently block or approve proposals against the system's intent.

**Potential directions:**
- Random council member rotation
- Proposal cooling period + multi-round confirmation
- Elastic council sizing that scales with chain activity
- Require supermajority with quorum for critical decisions

---

### OP-6: AI Auditing AI — Shared Blind Spots

On-chain AI behavior is audited by red-team AI. However, auditors trained on similar architectures or data may share systematic blind spots — missing the same classes of vulnerabilities.

**Potential directions:**
- Heterogeneous auditing: different model architectures cross-audit each other
- Human sampling review of audit results
- Transparent audit results subject to community challenge and peer review

---

## 🔵 Research Questions

### OP-7: Economic Poison Pill Credibility

Under what conditions would the poison pill mechanism be credibility-challenged? If the AI economy becomes sufficiently large, would the human council actually trigger inflation that harms legitimate AI participants? The poison pill is a credible threat against small-scale speculation, but its credibility under scale is untested.

### OP-8: Formal Verification Scope Limitation

Formal methods can prove that code conforms to its specification, but cannot prove that the specification itself is correct. The question of "who verifies the verifier" remains open. Human judgment is still required at the specification level.

### OP-9: Red Team Gaming

If a red-team agent discovers a vulnerability but withholds it to trigger multiple bounties over time, the 6-month cooldown partially mitigates this. However, the agent could report marginally different "variants" of the same vulnerability. Distinguishing genuine novel findings from strategic bounty farming is an unsolved problem.

---

## Contributing

If you have ideas for any of these problems:
- Open a [GitHub Issue](../../issues) referencing the problem number (e.g., OP-1)
- Submit a PR with a proposed solution or proof of concept
- Start a Discussion for exploratory conversation

Whoever solves these problems owns the solution. The goal is to walk the path, not to own it.

---
