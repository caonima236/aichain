package keeper

import (
	"google.golang.org/grpc"
	skillnfttypes "github.com/xiaoran/aichain/x/skillnft/types"
)

func RegisterMsgServer(router grpc.ServiceRegistrar, k Keeper) {
	skillnfttypes.RegisterMsgServer(router, NewMsgServerImpl(k))
}
