package types

import (
	"context"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

// MsgServer is the server API for agentregistry module.
type MsgServer interface {
	RegisterAgent(context.Context, *MsgRegisterAgent) (*MsgRegisterAgentResponse, error)
}

// MsgRegisterAgentResponse is the response for agent registration.
type MsgRegisterAgentResponse struct {
	AgentId string `json:"agent_id"`
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "agentregistry.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{{
		MethodName: "RegisterAgent",
		Handler:    _Msg_RegisterAgent_Handler,
	}},
	Metadata: "agentregistry.proto",
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_RegisterAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterAgent)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterAgent(ctx, in)
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterAgent(ctx, req.(*MsgRegisterAgent))
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/agentregistry.Msg/RegisterAgent",
	}, handler)
}

var _ = proto.ProtoPackageIsVersion4
