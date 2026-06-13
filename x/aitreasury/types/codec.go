package types

import (
	"context"
	"google.golang.org/grpc"
)

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "aitreasury.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{{
		MethodName: "CreateTreasuryProposal",
		Handler:    _Msg_CreateTreasuryProposal_Handler,
	}},
	Metadata: "aitreasury.proto",
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CreateTreasuryProposal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateTreasuryProposal)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { if srv == nil { return &MsgCreateTreasuryProposalResponse{}, nil }; return srv.(MsgServer).CreateTreasuryProposal(ctx, in) }
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateTreasuryProposal(ctx, req.(*MsgCreateTreasuryProposal))
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{FullMethod: "/aitreasury.Msg/CreateTreasuryProposal", Server: srv}, handler)
}
