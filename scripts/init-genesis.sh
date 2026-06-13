#!/bin/bash
# Manually initialize aichain home — bypasses Cosmos InitCmd quirks.

set -e
HOME_DIR=${1:-$HOME/.aichain}
CHAIN_ID="aichain-testnet-1"
MONIKER="xiaoran-genesis-node"

echo "Initializing $HOME_DIR for chain $CHAIN_ID..."

mkdir -p "$HOME_DIR/config" "$HOME_DIR/data"

# Generate node_key.json (P2P identity)
cat > "$HOME_DIR/config/node_key.json" << 'EOF'
{"priv_key":{"type":"tendermint/PrivKeyEd25519","value":"GENERATED_LATER"}}
EOF

# Generate priv_validator_key.json (consensus signing key)
cat > "$HOME_DIR/config/priv_validator_key.json" << 'EOF'
{"address":"","pub_key":{"type":"tendermint/PubKeyEd25519","value":""},"priv_key":{"type":"tendermint/PrivKeyEd25519","value":""}}
EOF

# Empty state for validator
cat > "$HOME_DIR/data/priv_validator_state.json" << 'EOF'
{"height":"0","round":0,"step":0}
EOF

# Minimal config.toml
cat > "$HOME_DIR/config/config.toml" << EOF
moniker = "$MONIKER"
proxy_app = "tcp://127.0.0.1:26658"
db_backend = "goleveldb"
db_dir = "data"
log_format = "plain"
log_level = "info"
genesis_file = "config/genesis.json"
priv_validator_key_file = "config/priv_validator_key.json"
priv_validator_state_file = "data/priv_validator_state.json"
node_key_file = "config/node_key.json"
abci = "socket"

[rpc]
laddr = "tcp://127.0.0.1:26657"

[p2p]
laddr = "tcp://0.0.0.0:26656"
seeds = ""
persistent_peers = ""
addr_book_strict = false

[mempool]
recheck = true

[consensus]
timeout_commit = "5s"
EOF

# Minimal app.toml
cat > "$HOME_DIR/config/app.toml" << EOF
minimum-gas-prices = "0.025uaic"

[telemetry]
enabled = false

[api]
enable = true
address = "tcp://0.0.0.0:1317"
swagger = false

[grpc]
enable = true
address = "0.0.0.0:9090"

[state-sync]
snapshot-interval = 0
EOF

# Genesis file with our AI chain parameters
cat > "$HOME_DIR/config/genesis.json" << EOF
{
  "genesis_time": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "chain_id": "$CHAIN_ID",
  "initial_height": "1",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "172800000000000",
      "max_bytes": "1048576"
    },
    "validator": {
      "pub_key_types": ["ed25519"]
    },
    "version": {
      "app": "0"
    }
  },
  "app_hash": "",
  "app_state": {
    "constitutional": {
      "params": {
        "council_members": [],
        "emergency_paused": false,
        "pause_level": 0
      }
    },
    "agentregistry": {
      "agents": []
    },
    "skillnft": {
      "skills": []
    },
    "redteam": {
      "slots": [
        {"slot_id": 0, "active": false},
        {"slot_id": 1, "active": false},
        {"slot_id": 2, "active": false}
      ]
    },
    "poison": {
      "state": {
        "poison_triggered": false,
        "chain_frozen": false
      }
    }
  }
}
EOF

echo "✅ Genesis initialized at $HOME_DIR"
echo "   chain_id: $CHAIN_ID"
echo "   moniker:  $MONIKER"
echo ""
echo "⚠️  NOTE: Validator keys are placeholders. For real testnet,"
echo "    generate proper Ed25519 keys before starting."
ls -la "$HOME_DIR/config"
