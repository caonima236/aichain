package types

import (
	"context"
	"google.golang.org/grpc"
)

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "agentdao.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{{
		MethodName: "CreateDAOProposal",
		Handler:    _Msg_CreateDAOProposal_Handler,
	}},
	Metadata: "agentdao.proto",
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CreateDAOProposal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateDAOProposal)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { if srv == nil { return &MsgCreateDAOProposalResponse{}, nil }; return srv.(MsgServer).CreateDAOProposal(ctx, in) }
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateDAOProposal(ctx, req.(*MsgCreateDAOProposal))
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{FullMethod: "/agentdao.Msg/CreateDAOProposal", Server: srv}, handler)
}
