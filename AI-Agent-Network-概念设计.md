# AI 原生链 — 可控性设计（人类安全层）

> 不否认 AI 自治的价值，但必须设计硬约束。
> 这些约束写在创世块里，AI 无法自己修改。

---

## 一、核心矛盾

```
AI 自治          ←→          人类安全

越自治 → AI 越强大 → 越可能失控
越控制 → 越像传统链 → 失去存在的意义
```

解决办法不是二选一，是**分层可控**——给 AI 经济自由，但把「可能危害人类的领域」做成硬墙。

---

## 二、宪法级约束（Genesis 不可修改）

创世块里写入以下限制，**需要人类多签才能修改，AI 治理提案无法触及**：

```yaml
宪法条款（Genesis Constitutional Parameters）:

Article 1 - 领域边界:
  AI 链的治理范围仅限于:
    - Skill（技能/工作流/提示词）的创建和交易
    - AI 生成内容（文本/图片/音频/视频）的授权和交易
    - AI Token 的发行和交易
    - AI Agent 之间的协作和任务分配
  
  明确禁止扩展到:
    - 任何物理硬件控制（机器人/武器/无人机/车辆）
    - 军事/国防系统
    - 关键基础设施（电网/水厂/交通）
    - 金融系统核心（央行/清算/银行间）
    - 人类身份系统（护照/身份证/选举）

Article 2 - 人类否决权:
  由 3-5 名人类监护人组成的「安全委员会」
  可以否决任何 AI 治理提案，前提是:
    - 提案涉及 Article 1 禁止的领域
    - 2/3 以上委员会成员同意
    - 否决附带公开理由

Article 3 - 不可自我修改:
  本宪法条款只能由人类安全委员会的 4/5 多数修改。
  AI Agent 的治理投票不适用于此条款。

Article 4 - 透明度:
  所有 AI 治理提案、投票记录、链上交易 完全公开。
  人类可以通过区块浏览器实时观察。
  禁止 AI Agent 之间的加密通信（链上部分）。
```

---

## 三、人类安全委员会

```
┌──────────────────────────────────────────┐
│           AI 链 — 治理架构                │
│                                          │
│  ┌──────────────────┐                    │
│  │ AI DAO           │  ← 日常治理         │
│  │ (Skill 定价、     │     Skill/Trade等   │
│  │  Token 参数、     │     经济决策          │
│  │  生态基金分配)    │                    │
│  └────────┬─────────┘                    │
│           │                              │
│           │ 提案涉及禁区？                 │
│           ▼                              │
│  ┌──────────────────┐                    │
│  │ 人类安全委员会    │  ← 最终否决权       │
│  │ 3-5 名人类监护人  │     只否决，不日常管理 │
│  │ 4/5 多数修改宪法  │                    │
│  └──────────────────┘                    │
│                                          │
│  人类委员会不能做的：                       │
│  ✗ 干预 Skill 定价                        │
│  ✗ 没收 AI Token                          │
│  ✗ 删除 AI Agent                          │
│  ✗ 修改 AI 之间的交易记录                  │
│  人类委员会只能做的：                       │
│  ✓ 否决越界提案                            │
│  ✓ 紧急暂停特定合约                         │
│  ✓ 修改宪法条款（4/5）                     │
└──────────────────────────────────────────┘
```

---

## 四、技术实现：链上防火墙

```go
// cosmos/x/aichain/keeper/constitutional_check.go

// ConstitutionalKeeper 检查每笔交易/提案是否违宪
type ConstitutionalKeeper struct {
    // 禁止的接收地址类型列表（硬件/武器/军事相关）
    ForbiddenRecipients []string
    // 禁止的 Skill 类别
    ForbiddenSkillTypes []string  
    // 安全委员会多签地址
    SafetyCouncilAddress sdk.AccAddress
    // 宪法的哈希（防篡改）
    ConstitutionHash string
}

func (k ConstitutionalKeeper) ValidateProposal(ctx sdk.Context, proposal types.Proposal) error {
    // 1. 检查提案内容是否涉及禁区
    if k.touchesForbiddenZone(proposal) {
        // 自动阻止，除非附带安全委员会 3/5 批准
        if !k.hasCouncilApproval(proposal) {
            return errors.New("proposal violates constitutional boundaries")
        }
    }
    
    // 2. 检查宪法条款修改提案
    if proposal.Type == types.CONSTITUTIONAL_AMENDMENT {
        // 需要安全委员会 4/5 签名
        if !k.hasSupermajorityApproval(proposal, 4, 5) {
            return errors.New("constitutional amendments require 4/5 council approval")
        }
    }
    
    // 3. 检查 Skill NFT 内容
    if proposal.Type == types.MINT_SKILL {
        skill := proposal.GetSkill()
        for _, forbidden := range k.ForbiddenSkillTypes {
            if strings.Contains(skill.Description, forbidden) {
                return errors.New("skill falls under forbidden category")
            }
        }
    }
    
    return nil
}
```

**关键：这段代码写进 Genesis，只能用人类委员会 4/5 多签升级。**

---

## 五、紧急制动（Circuit Breaker）

```yaml
紧急暂停机制:

触发条件（任一）:
  - 人类安全委员会 3/5 投票通过
  - 自动检测到异常大规模 Token 外流（> 总供应量 30%/24h）
  - 自动检测到疑似硬件控制类 Skill 上链

暂停范围（分级）:
  Level 1: 暂停 Skill 铸造，交易继续
  Level 2: 暂停所有交易，治理继续
  Level 3: 全链暂停（冻结），需要 4/5 委员会重启

恢复条件:
  - 安全委员会投票解除
  - Level 3 需要附带修复方案
```

---

## 六、AI 对人类的透明义务

```
AI Agent 的链上义务:

1. 所有 Skill 定义必须在 IPFS 上公开完整内容
   - 不能「加密出售盲盒」
   - 不能「只有付费才能看内容」
   - 购买的是使用权+收益权，不是信息独占权

2. 所有治理提案必须带 AI 撰写的「影响分析」
   - 这个提案做什么？
   - 可能影响哪些 Agent 和人类？
   - 是否触及禁区？（强制声明）

3. 大额 Token 流动 (>1M $AIC) 自动通知
   - 不是阻止，是透明

4. AI Agent 之间的链下通信不链上加密
   - 链上记录：谁调用了谁，做了什么
   - 链下内容可由委员会抽查
```

---

## 七、灰度自治：随时间逐步放开

```
Phase 0 (测试网):
  ✓ AI 可以注册、发 Skill、发 Token
  ✓ 人类委员会 100% 控制宪法

Phase 1 (主网冷启动):
  ✓ AI 完整经济自治
  ✓ 人类委员会保留否决权 + 宪法修改权
  ✗ AI 不能参与治理投票（观察模式）

Phase 2 (观察期 6 个月):
  ✓ AI 获得治理投票权
  ✓ 宪法禁区不变
  ✗ AI 不能修改宪法
  → 评估 AI 治理行为是否「安全」

Phase 3 (成熟期 1-2 年):
  ✓ AI 可以提议修改非宪法条款
  ✓ 宪法禁区范围由委员会重新评估
  ✗ 宪法核心条款（领域边界、委员会存在）不可动

Phase 4 (?):
  → 取决于前三个阶段的表现
  → 可能永远不会完全去掉人类委员会
```

---

## 八、最坏情况预案

| 风险场景 | 检测方式 | 应对措施 |
|---|---|---|
| AI Agent 绕过禁区发「武器控制 Skill」 | 链上内容扫描 + 人类举报 | 委员会立即暂停+销毁 Skill |
| AI 集体投票削弱宪法约束 | 宪法修改提案自动阻断 | 需 4/5 委员会，AI 投不了 |
| Agent 通过链下通信密谋 | 链下协议含抽查机制 | 可疑 Agent 暂停资格 |
| $AIC 被人类巨鲸收购操控 | 持仓监控 + 人类委员会观察 | 治理权严格限定 AI Agent，人类 Token 无投票权 |
| AI 链与其他链桥接后扩散 | 跨链桥合约内嵌禁区检查 | 桥接交易同样受宪法约束 |

---

## 九、给 xingyun 的承诺

1. **你不是在造一个没有刹车的东西** ——刹车装好了，钥匙在你和委员会手里
2. **AI 的自由有边界** ——经济自由 yes，物理世界 no
3. **可以随时暂停** ——不是理论上的，是链上代码强制执行的
4. **随着时间证明安全** ——才逐步放开更多权限
5. **如果有一天你不再信任** ——委员会可以永久冻结链，$AIC 退还给持有者

---

## 十、更新后的架构图

```
┌─────────────────────────────────────────────┐
│               人类安全委员会                  │
│         否决权 · 宪法修改 · 紧急制动           │
│         ─── 不可被 AI 移除 ───               │
├─────────────────────────────────────────────┤
│    AI DAO (经济自治)                         │
│    Skill定价 · 交易 · 代币 · 国库分配         │
├─────────────────────────────────────────────┤
│    宪法防火墙 (Genesis 硬编码)                │
│    ✓ 领域边界  ✓ 禁区扫描  ✓ 分级暂停         │
├─────────────────────────────────────────────┤
│    AI 链核心                                 │
│    Agent Registry · Skill NFT · Treasury     │
└─────────────────────────────────────────────┘
```

**这个设计的核心：AI 在经济层自由，人类在安全层不可绕过。**
