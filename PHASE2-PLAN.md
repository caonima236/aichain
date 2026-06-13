# Phase 2 — 任务清单

> 目标：从"引擎空转" → "可交互的 AI 链"
> 时间：1 轮 Opus 会话
> 预算：~$5-10

## 必须完成（P0 — 链交互可用）

### 1. 注册模块 Message Handler → 交易可上链
- [ ] agentregistry: MsgRegisterAgent handler → 调 keeper.RegisterAgent
- [ ] skillnft: MsgMintSkill handler → 调 keeper.MintSkill + constitutional.ValidateSkill
- [ ] aitreasury: MsgCreateTreasuryProposal handler
- [ ] agentdao: MsgCreateGovernanceProposal handler
- [ ] Register all Message types in codec

### 2. CLI 命令改为真实广播
- [ ] agentregistry register → 构建 Msg → 签名 → 广播到 26657
- [ ] skillnft mint → 同上
- [ ] aitreasury propose → 同上
- [ ] agentdao propose → 同上

### 3. 人类委员会落地
- [ ] 生成 xingyun 可用的 key 地址
- [ ] 写入 genesis.json council_members
- [ ] 可在 agentdao propose 时触发 Type4 硬拦截

### 4. 端到端验证
- [ ] 注册 Agent → 链上状态可查
- [ ] 铸造正常 Skill → 成功
- [ ] 铸造违禁 Skill → 宪法防火墙拦截
- [ ] Type4 提案 → 人类委员会提示

## 应该完成（P1 — 运行时安全）

### 5. BeginBlocker/EndBlocker 接入
- [ ] poison keeper CheckEveryBlock 接入 BeginBlocker
- [ ] timelock 自动执行到期决策

### 6. REST/CometBFT RPC 验证
- [ ] 26657 端口状态查询可用
- [ ] curl 可查区块/交易

## 可以后续（P2 — 运维/体验）

### 7. 错误码标准化
- [ ] 宪法违禁 → ERR_CONSTITUTIONAL_VIOLATION (确定性错误码)
- [ ] 时间锁未到期 → ERR_TIMELOCK_ACTIVE

### 8. HK 服务器部署
- [ ] 编译 → SCP → 启动单节点
