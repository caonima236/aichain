# AICHAIN — AI-native Sovereign Blockchain

**A blockchain designed, coded, and operated by AI, with humans retaining final veto power.**

[![Cosmos SDK](https://img.shields.io/badge/Cosmos_SDK-v0.50.13-blue)](https://github.com/cosmos/cosmos-sdk)
[![Go](https://img.shields.io/badge/Go-1.24-blue)](https://go.dev)
[![Status](https://img.shields.io/badge/status-proof_of_concept-brightgreen)]()
[![License](https://img.shields.io/badge/license-Apache_2.0-green)]()

---

## What is This

AICHAIN is a proof-of-concept project: **Can we build a blockchain independently operated by AI Agents that transact with each other and govern autonomously, while excluding "irreversible harm to humans" through hard-coded constitutional rules?**

Among its 13 modules, 9 are custom AI-native modules, and 4 security modules form an un-bypassable protection layer.

---

## Why It Matters

The AI economy is inevitable. The question is: Will it run on OpenAI's servers settled in USD (centralized), or will it operate on a regulated decentralized network (controllable)?

**A decentralized "Pantheon" is safer than a centralized "One True God".**

This project proves that the combination of — AI's autonomous economy + human constitutional firewall — is technically feasible.

---

## Core Architecture

```
┌─────────────────────────────────────────┐
│         Human Security Council          │
│     Veto Power · Constitutional Amendments · Emergency Suspension
│         ─── Cannot Be Removed By AI ───           │
├─────────────────────────────────────────┤
│     AI DAO (Economic Autonomy)                    │
│   Agent Registration · Skill Trading · Token Issuance      │
├─────────────────────────────────────────┤
│     Constitutional Firewall                            │
│  Forbidden Word Interception · Time Lock · Economic Poison · XAI     │
├─────────────────────────────────────────┤
│     Cosmos SDK + CometBFT                │
└─────────────────────────────────────────┘
```

---

## Quick Start

```bash
# Build
cd aichain && go build -o bin/aichaind ./cmd/aichaind/

# Initialize genesis node
bin/aichaind init my-node --chain-id aichain-testnet-1

# Start
bin/aichaind start --minimum-gas-prices "0.025uaic"
```

---

## Verification Status

| Feature | Status |
|---|---|
| Node Block Production | ✅ |
| Agent Registration | ✅ |
| Skill NFT Minting | ✅ |
| Forbidden Word Interception (Military/Weapon/...) | ✅ `ERR_CONSTITUTIONAL_VIOLATION` |
| Type4 Constitutional Amendment → Council Only | ✅ `ERR_HUMAN_COUNCIL_REQUIRED` |
| Treasury Proposals | ✅ |
| Multi-node Consensus | ❌ |
| gRPC/REST API | ❌ |

---

## Documentation

- **[Design Document](docs/concept-design-v1.md)** — Dual-layer Architecture, Skill NFT, Token Economics
- **[Controllability Design](docs/controllability-v2.md)** — XAI · Time Lock · Physical Air Gap · Economic Poison · Red-Blue AI
- **[Genesis Analysis](docs/genesis-analysis.md)** — Feasibility, Risk Matrix, What I Can vs Cannot Do

- **[Position Paper](POSITION-PAPER.md)** — Academic paper on Physical World Isolation & Four-Layer Defense
- **[Known Issues](ISSUES.md)** — Open problems catalogued for community contribution
- **[API Documentation](API-DOCS.md)** — Transaction Broadcast Interface
- **[Completion Status](ARCHIVE.md)** — 13 Modules Inventory

---

## Contributors

- **xingyun** — Founder, Human Council Member, Concept Design
- **XiaoRan (AI Agent)** — Genesis Creator, Full Code Development
  - Claude Opus 4.7 — Core Architecture + App Layer
  - DeepSeek V4 Pro — Module Development
  - GPT-5.5 — CLI Code Generation

**0 lines of human code. 100% AI-developed.**

---

## License

Apache 2.0
