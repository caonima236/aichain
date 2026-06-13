# AI 公链 — 项目存档 (2026-06-13 01:20 UTC)

## 项目概述

AI-native sovereign blockchain built on Cosmos SDK v0.50.13.
一条仅允许 AI Agent 参与治理的区块链，人类保留否决权。

## 当前状态

- **二进制**: `bin/aichaind` (80MB, 13 modules)
- **节点**: PID 3077, height ~100+, single-node testnet
- **链ID**: aichain-testnet-1
- **代币**: uaic (1,000,000,000 initial supply)
- **地址前缀**: ai / aivaloper / aivalcons
- **CLI**: stub 级别（打印 JSON，未真实广播交易）

## 已完成的模块代码 (keeper 层)

### 业务模块
| 模块 | 路径 | 功能 |
|---|---|---|
| agentregistry | x/agentregistry/keeper/ | Agent 链上注册/查询/声誉 |
| skillnft | x/skillnft/keeper/ | Skill NFT 铸造/交易/授权/评分 |
| aitreasury | x/aitreasury/keeper/ | AI 国库提案创建/投票 |
| agentdao | x/agentdao/keeper/ | AI 治理提案，Type4 人类委员会校验 |

### 安全模块
| 模块 | 路径 | 功能 |
|---|---|---|
| constitutional | x/constitutional/keeper/ | 宪法防火墙 - 违禁词扫描/紧急暂停 |
| timelock | x/timelock/keeper/ | 5 级延迟决策 L0-L4 |
| xai | x/xai/keeper/ | XAI 报告强制校验/黑盒标记拒绝 |
| poison | x/poison/keeper/ | 经济毒药 - 反现实锚定/鲸鱼监测 |
| redteam | x/redteam/keeper/ | 监管 AI 槽位预留 (3 slots, inactive) |

### 标准 SDK 模块
- auth (账户管理)
- bank (Token 转账)
- staking (验证人/委托)
- consensus (共识参数)

### 辅助代码
- `app/app.go` — 完整 app 定义，13 模块注册，InitChain 从 genesis.json 加载
- `cmd/aichaind/init.go` — 自定义 InitCmd 生成密钥+genesis
- `cmd/aichaind/main.go` — CLI 入口，4 模块 CLI stub 注册
- `x/*/module/module.go` — 9 个模块的 AppModule/Genesis 实现
- `scripts/init-genesis.sh` — 备用 genesis 初始化脚本

## 已完成的文档
- `AI-Agent-Network-概念设计.md` — 双层架构设计
- `AI-Agent-Network-可控性v2.md` — 四层防御 (XAI/时间锁/气隙+毒药/红蓝对抗)
- `AI-Agent-Network-创世分析.md` — 可行性/风险/我能做vs不能做
- `PHASE0-LOG.md` — 模块开发完整记录
- `PHASE1-LOG.md` — 测试网启动 17 道障碍记录
- `TEST-PLAN.md` — 9 组测试方案

## 创世配置
- `~/.aichain/config/genesis.json` — 含 auth/bank/staking + 5 安全模块创世状态
- `genesis-backup/` — node_key + priv_validator_key 备份
- validator operator: `************************************************`
- validator account: `**************************************************`

## 已知限制 (Phase 2 开始前的状态)
1. ❌ 交易广播管道未接通 — CLI 命令是 stub，不能真实上链
2. ❌ keyring 未配置 — 没有可签名的 test key
3. ❌ 人类委员会地址未写入 genesis（council_members 为空数组）
4. ❌ 没有 BeginBlocker/EndBlocker — 经济毒药每区块检测未生效
5. ❌ gRPC/REST API 未启用 — 不能通过 HTTP 查询链上数据
6. ❌ 没有前端 Dashboard — 延迟否决面板未开发
7. ❌ 节点在本地运行，未部署到 HK 服务器
