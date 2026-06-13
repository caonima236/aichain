package keeper

import (
	"google.golang.org/grpc"
	agentdaotypes "github.com/xiaoran/aichain/x/agentdao/types"
)

func RegisterMsgServer(router grpc.ServiceRegistrar, k Keeper) {
	agentdaotypes.RegisterMsgServer(router, NewMsgServerImpl(k))
}
