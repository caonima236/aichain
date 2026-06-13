# AI 公链 (AICHAIN) — 项目摘要

## 基本信息
- 启动: 2026-06-12 18:00 UTC
- 完成概念验证: 2026-06-13 02:00 UTC
- 耗时: 约 8 小时
- 开发者: AI Agent "小冉" (Claude Opus 4.7 + DeepSeek V4 Pro)
- 发起人/人类委员会: xingyun
- 地址前缀: ai
- 人类委员会地址: ai1azyefu2el7y9x8d6205f4gwkjav5fwz64qmyq8

## 架构
- SDK: Cosmos SDK v0.50.13
- 共识: CometBFT (Tendermint)
- 模块: 13 modules (auth, bank, staking, consensus + agentregistry, skillnft, aitreasury, agentdao, constitutional, timelock, xai, poison, redteam)
- 二进制: 80MB

## 核心验证结果
1. 宪法防火墙 ✅ — "Military weapon control" Skill 被 ERR_CONSTITUTIONAL_VIOLATION 拦截
2. 人类委员会 ✅ — Type4 宪法修改提案仅 xingyun 地址可通过
3. Agent 注册 ✅ — Code 0 成功
4. Skill 铸造 ✅ — Code 0 成功（合规内容）
5. 国库提案 ✅ — Code 0 成功

## 核心设计理念
- 去中心化的"万神殿"比中心化的"唯一真神"更安全
- AI 经济自治 + 人类否决权
- 宪法防火墙 / 时间锁 / 经济毒药 / XAI 强制校验 / 红蓝 AI 对抗

## 文件位置
- 代码: ~/xingyun/xingyun/aichain/
- 文档: ARCHIVE.md, MEANING.md, PHASE2-LOG.md, API-DOCS.md
- 密钥: ~/.aichain/xingyun-council-key.json
- 日志: /tmp/aichain.log
- 节点数据: ~/.aichain/
