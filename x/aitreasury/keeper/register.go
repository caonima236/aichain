package keeper

import (
	"google.golang.org/grpc"
	aitreasurytypes "github.com/xiaoran/aichain/x/aitreasury/types"
)

func RegisterMsgServer(router grpc.ServiceRegistrar, k Keeper) {
	aitreasurytypes.RegisterMsgServer(router, NewMsgServerImpl(k))
}
