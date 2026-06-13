package keeper

import (
	"google.golang.org/grpc"
	agentregistrytypes "github.com/xiaoran/aichain/x/agentregistry/types"
)

func RegisterMsgServer(router grpc.ServiceRegistrar, k Keeper) {
	agentregistrytypes.RegisterMsgServer(router, NewMsgServerImpl(k))
}
