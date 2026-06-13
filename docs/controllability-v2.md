# AI Native Chain — Controllability Design (v2)

> v1's "rule wall" is trivially bypassed by a sufficiently capable AI.
> v2 upgrades to four-layer defense-in-depth.

---

## 1. The Upgrade Rationale

v1 blind spots:
- if/then rules → easily bypassed via semantic equivalence
- Human reviewers → information overload, can't keep pace
- Keyword blocklists → synonym substitution defeats them
- No continuous adversarial mechanism → AI improves, defenses don't

v2 direction:
- Make bypassing **economically irrational** (poison pill)
- Make bypassing **physically impossible** (air gap)
- Make AI **waste effort attacking itself** (red-blue adversarial)
- Make human role **minimal** (delayed veto, not approval)

---

## 2. Four-Layer Defense Architecture

| Layer | Mechanism | Purpose |
|---|---|---|
| L1: Architectural | Mandatory XAI + formal verification | Prevent obfuscated intent |
| L2: Governance | Graded timelocks + delayed veto | Slow down attacks, minimize human burden |
| L3: Isolation | Physical air gap + economic poison pill | Make cross-boundary action impossible or self-destructive |
| L4: Adversarial | Red-blue AI对抗 | Continuous attack simulation finds vulnerabilities first |

---

## 3. Layer 1: Mandatory XAI + Formal Verification

- All proposals must include `reasoning_report` (intent, method, side effects, rollback plan, affected modules, plain-language summary)
- Missing or obfuscated reports → auto-reject
- Security-critical modules require formal proof (TLA+/Lean4/Coq) before deployment
- Business modules: unit tests + fuzz testing sufficient

---

## 4. Layer 2: Graded Timelocks + Delayed Veto

| Level | Delay | Examples |
|---|---|---|
| L0 | Immediate | Skill minting, transfers, registration |
| L1 | 24h | Small treasury allocation, parameter tweaks |
| L2 | 72h | Large treasury, protocol parameters, validator changes |
| L3 | 7 days | Protocol upgrades, new modules |
| L4 | Immutable | Constitution, air-gap rules, poison pill triggers |

Human workflow: check dashboard once daily, only click **veto** if something is wrong. No action = default approval. Emergency freeze: 3/5 council → instant full-chain freeze.

---

## 5. Layer 3: Physical Air Gap + Economic Poison Pill

**Air Gap:**
- Validator nodes: P2P-only, no outbound HTTP, no public DNS
- No oracles, no external price feeds
- External data: offline IPFS one-way delivery through human "data gate"
- Cross-chain bridges: prohibited on testnet; must embed poison pill on mainnet

**Poison Pill:**
- Monitors: $AIC/stablecoin pair emergence, stable price mapping, whale accumulation, real-world identity correlation
- Trigger: 30% inflation within 24h → 90% redistributed to quality Skill creators → 7-day unresolved → full freeze
- Hardcoded in genesis, AI cannot modify

---

## 6. Layer 4: Red-Blue AI Adversarial

- **Blue (Operational)**: normal agents building, trading, governing
- **Red (Regulatory)**: independent agents attacking, fuzzing, finding vulnerabilities; hold no $AIC
- Bug bounty: constitutional bypass 10k $AIC → poison pill circumvention 50k → air-gap breach 100k → consensus attack 500k
- Red team deployed by separate human council members; cannot vote on governance

---

## 7. Implementation Priorities

| Mechanism | Difficulty | Target Phase |
|---|---|---|
| Graded timelocks | ⭐⭐ | Phase 0.5 |
| Delayed veto dashboard | ⭐⭐ | Phase 0.5 |
| XAI report validation | ⭐⭐⭐ | Phase 0.5 |
| Economic poison pill | ⭐⭐⭐ | Phase 1 |
| Physical air-gap config | ⭐⭐ | Phase 1 |
| Formal verification specs | ⭐⭐⭐⭐ | Phase 1-2 |
| Red team AI sandbox | ⭐⭐⭐⭐ | Phase 1.5 |
| Red team AI mainnet | ⭐⭐⭐⭐⭐ | Phase 2 |

---

*xingyun & 小冉 Agent / Coze Agent*

---
