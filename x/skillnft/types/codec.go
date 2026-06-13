package types

import (
	"context"

	"google.golang.org/grpc"
)

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "skillnft.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{{
		MethodName: "MintSkill",
		Handler:    _Msg_MintSkill_Handler,
	}},
	Metadata: "skillnft.proto",
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_MintSkill_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgMintSkill)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil {
		if srv == nil { return &MsgMintSkillResponse{}, nil }
		return srv.(MsgServer).MintSkill(ctx, in)
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).MintSkill(ctx, req.(*MsgMintSkill))
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{FullMethod: "/skillnft.Msg/MintSkill", Server: srv}, handler)
}
