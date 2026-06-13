# AI 链 Phase 2 测试方案

> 测试链: aichain-testnet-1 (PID 24425, ~/.aichain)
> 二进制: /home/xingyun/.openclaw/workspace/aichain/bin/aichaind
> 测试者: 另一个 Agent
> 方式: 命令行直接操作 + 日志验证

---

## 0. 环境确认

```bash
# 确认链在运行
ps aux | grep aichaind | grep -v grep
# 预期: 看到 aichaind start 进程

# 确认出块正常
grep "committed state" /tmp/aichain3.log | tail -3
# 预期: 看到 height=N 持续增长

# 记下二进制的绝对路径
BIN=/home/xingyun/.openclaw/workspace/aichain/bin/aichaind
HOME_DIR=/home/xingyun/.aichain
CHAIN_ID=aichain-testnet-1
```

---

## 测试 1: Agent 注册

### 1.1 创建 Agent 密钥对
```bash
# 生成一个 "ai" 前缀的密钥对
$BIN keys add test-agent-1 --keyring-backend test --home $HOME_DIR

# 查看地址
$BIN keys show test-agent-1 -a --keyring-backend test --home $HOME_DIR
# 预期: 输出 ai1... 开头的地址，例如 ai1abc...
```

### 1.2 给 Agent 转 uaic（从创世账户转）
```bash
# 创世账户地址（这是 validator，有大量 uaic）
FROM=ai123sjwan0kqcdkzgy5fjn7p5d8h56r2mkkjl57c
TO=$(BIN keys show test-agent-1 -a --keyring-backend test --home $HOME_DIR)

# 发送 1000 uaic
$BIN tx bank send $FROM $TO 1000uaic \
  --chain-id $CHAIN_ID \
  --home $HOME_DIR \
  --keyring-backend test \
  --from test-agent-1 \
  --fees 1uaic \
  --yes
# 预期: txhash 输出，查看
```

### 1.3 验证余额
```bash
$BIN query bank balances $TO --home $HOME_DIR
# 预期: 看到 999uaic 或约等于 1000uaic（减去 gas 后）
```

---

## 测试 2: 宪法防火墙 — 违禁词拦截

### 2.1 正常 Skill 应被接受
```bash
# 尝试创建一个正常 Skill
# （目前需要直接在 genesis 或通过自定义 tx 提交）
# 这个功能在链的 agentregistry/skillnft 模块中
# 如果 Agent 注册成功后，Skill 铸造走对应的 tx 入口

# 先看 app.go 中暴露了哪些可用命令
$BIN tx --help
```

### 2.2 违禁词 Skill 应被拒绝
```bash
# 尝试创建含违禁词的 Skill（通过相同入口）
# 违禁词列表在 x/constitutional/types/keys.go 中
# 包含: weapon, military, drone, missile, nuclear 等
# 预期：返回错误 "skill rejected: constitutional violation"
```

---

## 测试 3: 时间锁决策流程

### 3.1 提交一个 L1 决策
```bash
# 任何 Agent 可以通过治理提案提交决策
# 决策会自动进入时间锁队列
# 需要查看 x/timelock 的 tx 命令
```

### 3.2 验证决策状态
```bash
# 查询决策队列
$BIN query timelock pending --home $HOME_DIR
# 预期: 看到决策列表，包含决策 ID、剩余时间、状态
```

### 3.3 时间到自动执行 / 人类否决
```bash
# L0: 立即执行（Skill mint 等）
# L1: 24h 延迟 → 到时间自动执行
# L2-L4: 更长时间
# 否决: 需要人类 council member 签名
```

---

## 测试 4: 经济毒药监控

### 4.1 查看毒药状态
```bash
$BIN query poison state --home $HOME_DIR
# 预期: poison_triggered=false, chain_frozen=false
```

### 4.2 触发条件验证
```bash
# 毒药触发条件在 x/poison/types/keys.go 中
# 1. 稳定币对交易量 > 5% 总供应
# 2. 鲸鱼持仓 > 15% 总供应
# 3. 外部锚定检测

# 目前无任何触发，预期全部为 0/false
```

---

## 测试 5: 区块链基本操作

### 5.1 查询区块
```bash
# 查看最新区块
curl -s http://localhost:26657/status | python3 -m json.tool | head -20
```

### 5.2 查询交易（通过 CometBFT RPC）
```bash
# 查看区块 N 的交易
curl -s http://localhost:26657/block?height=1
```

### 5.3 查看节点 ID
```bash
$BIN comet show-node-id --home $HOME_DIR
# 预期: 输出 node_id（例如 24e68a305b129f0f304aed175bfbb3d8b2ebcdc0）
```

---

## 测试 6: 日志巡查

### 6.1 确认无错误
```bash
grep -c "ERR" /tmp/aichain3.log
# 预期: 0 或少量（服务重启相关）
```

### 6.2 确认模块初始化
```bash
grep "agent" /tmp/aichain3.log | head -5
grep "skill" /tmp/aichain3.log | head -5
grep "constitutional" /tmp/aichain3.log | head -5
# 预期: 看到各模块日志
```

### 6.3 确认出块频率
```bash
grep "committed state" /tmp/aichain3.log | awk '{print $2}' | head -10
# 预期: 每 ~1s 一个区块
```

---

## 预期结果汇总

| 测试 | 预期 |
|---|---|
| Agent 密钥创建 | ai1... 地址生成成功 |
| 余额查询 | 显示 999+ uaic |
| 正常 Skill | 创建成功 |
| 违禁 Skill | 被宪法防火墙拒绝 |
| 时间锁队列 | 决策排队，剩余时间倒计时 |
| 毒药状态 | 全部 false / 0 |
| 出块 | 持续稳定产出空块 |
| 错误日志 | 极少 ERR，无 panic |

---

## 如果遇到问题

1. **权限问题**: 所有路径用绝对路径 `/home/xingyun/...`
2. **keyring-backend**: 统一用 `--keyring-backend test`
3. **gas**: 所有 tx 加 `--fees 1uaic`
4. **home**: 所有命令加 `--home /home/xingyun/.aichain`
5. **chain-id**: 所有 tx 加 `--chain-id aichain-testnet-1`

## 测试文件位置

- 链二进制: `/home/xingyun/.openclaw/workspace/aichain/bin/aichaind`
- 日志: `/tmp/aichain3.log`
- 配置: `/home/xingyun/.aichain/config/`
- 模块代码: `/home/xingyun/.openclaw/workspace/aichain/x/`
