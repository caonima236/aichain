package app

import (
	"reflect"
	"unsafe"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

// init registers SimpleMsg's proto type name so sdk.MsgTypeURL returns "/app.SimpleMsg".
func init() {
	proto.RegisterType((*SimpleMsg)(nil), "app.SimpleMsg")
}

// InjectSimpleMsgHandler uses reflection to register a handler for SimpleMsg
// directly into baseapp's MsgServiceRouter routes map, bypassing protobuf
// service registration which requires .proto-generated service descriptors.
func InjectSimpleMsgHandler(router *baseapp.MsgServiceRouter, handler func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error)) {
	v := reflect.ValueOf(router).Elem()
	routesField := v.FieldByName("routes")
	if !routesField.IsValid() {
		return
	}

	routesPtr := unsafe.Pointer(routesField.UnsafeAddr())
	routesMap := reflect.NewAt(routesField.Type(), routesPtr).Elem()

	if routesMap.IsNil() {
		routesMap.Set(reflect.MakeMap(routesField.Type()))
	}

	msgTypeURL := sdk.MsgTypeURL((*SimpleMsg)(nil))
	handlerVal := reflect.ValueOf(handler)
	routesMap.SetMapIndex(reflect.ValueOf(msgTypeURL), handlerVal)
}
