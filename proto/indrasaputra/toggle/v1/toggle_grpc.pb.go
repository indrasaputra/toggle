// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package togglev1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ToggleCommandServiceClient is the client API for ToggleCommandService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ToggleCommandServiceClient interface {
	// Create a new toggle.
	//
	// This endpoint creates a new toggle with provided key and description.
	// The description can be left empty, but the key must exists.
	// The key must be unique and it can only contain alphanumeric and dash.
	// The key will be converted to lower case.
	CreateToggle(ctx context.Context, in *CreateToggleRequest, opts ...grpc.CallOption) (*CreateToggleResponse, error)
	// Enable a toggle.
	//
	// This endpoint set toggle's usability to active.
	// Its *isEnabled* attribute will be set to true.
	EnableToggle(ctx context.Context, in *EnableToggleRequest, opts ...grpc.CallOption) (*EnableToggleResponse, error)
	// Disable a toggle.
	//
	// This endpoint set toggle's usability to inactive.
	// Its *isEnabled* attribute will be set to false.
	DisableToggle(ctx context.Context, in *DisableToggleRequest, opts ...grpc.CallOption) (*DisableToggleResponse, error)
	// Delete a toggle.
	//
	// This endpoint deletes a toggle by its key.
	// The operation is hard-delete, thus the toggle will be gone forever.
	DeleteToggle(ctx context.Context, in *DeleteToggleRequest, opts ...grpc.CallOption) (*DeleteToggleResponse, error)
}

type toggleCommandServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewToggleCommandServiceClient(cc grpc.ClientConnInterface) ToggleCommandServiceClient {
	return &toggleCommandServiceClient{cc}
}

func (c *toggleCommandServiceClient) CreateToggle(ctx context.Context, in *CreateToggleRequest, opts ...grpc.CallOption) (*CreateToggleResponse, error) {
	out := new(CreateToggleResponse)
	err := c.cc.Invoke(ctx, "/proto.indrasaputra.toggle.v1.ToggleCommandService/CreateToggle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toggleCommandServiceClient) EnableToggle(ctx context.Context, in *EnableToggleRequest, opts ...grpc.CallOption) (*EnableToggleResponse, error) {
	out := new(EnableToggleResponse)
	err := c.cc.Invoke(ctx, "/proto.indrasaputra.toggle.v1.ToggleCommandService/EnableToggle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toggleCommandServiceClient) DisableToggle(ctx context.Context, in *DisableToggleRequest, opts ...grpc.CallOption) (*DisableToggleResponse, error) {
	out := new(DisableToggleResponse)
	err := c.cc.Invoke(ctx, "/proto.indrasaputra.toggle.v1.ToggleCommandService/DisableToggle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toggleCommandServiceClient) DeleteToggle(ctx context.Context, in *DeleteToggleRequest, opts ...grpc.CallOption) (*DeleteToggleResponse, error) {
	out := new(DeleteToggleResponse)
	err := c.cc.Invoke(ctx, "/proto.indrasaputra.toggle.v1.ToggleCommandService/DeleteToggle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ToggleCommandServiceServer is the server API for ToggleCommandService service.
// All implementations must embed UnimplementedToggleCommandServiceServer
// for forward compatibility
type ToggleCommandServiceServer interface {
	// Create a new toggle.
	//
	// This endpoint creates a new toggle with provided key and description.
	// The description can be left empty, but the key must exists.
	// The key must be unique and it can only contain alphanumeric and dash.
	// The key will be converted to lower case.
	CreateToggle(context.Context, *CreateToggleRequest) (*CreateToggleResponse, error)
	// Enable a toggle.
	//
	// This endpoint set toggle's usability to active.
	// Its *isEnabled* attribute will be set to true.
	EnableToggle(context.Context, *EnableToggleRequest) (*EnableToggleResponse, error)
	// Disable a toggle.
	//
	// This endpoint set toggle's usability to inactive.
	// Its *isEnabled* attribute will be set to false.
	DisableToggle(context.Context, *DisableToggleRequest) (*DisableToggleResponse, error)
	// Delete a toggle.
	//
	// This endpoint deletes a toggle by its key.
	// The operation is hard-delete, thus the toggle will be gone forever.
	DeleteToggle(context.Context, *DeleteToggleRequest) (*DeleteToggleResponse, error)
	mustEmbedUnimplementedToggleCommandServiceServer()
}

// UnimplementedToggleCommandServiceServer must be embedded to have forward compatible implementations.
type UnimplementedToggleCommandServiceServer struct {
}

func (UnimplementedToggleCommandServiceServer) CreateToggle(context.Context, *CreateToggleRequest) (*CreateToggleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateToggle not implemented")
}
func (UnimplementedToggleCommandServiceServer) EnableToggle(context.Context, *EnableToggleRequest) (*EnableToggleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableToggle not implemented")
}
func (UnimplementedToggleCommandServiceServer) DisableToggle(context.Context, *DisableToggleRequest) (*DisableToggleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableToggle not implemented")
}
func (UnimplementedToggleCommandServiceServer) DeleteToggle(context.Context, *DeleteToggleRequest) (*DeleteToggleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteToggle not implemented")
}
func (UnimplementedToggleCommandServiceServer) mustEmbedUnimplementedToggleCommandServiceServer() {}

// UnsafeToggleCommandServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ToggleCommandServiceServer will
// result in compilation errors.
type UnsafeToggleCommandServiceServer interface {
	mustEmbedUnimplementedToggleCommandServiceServer()
}

func RegisterToggleCommandServiceServer(s grpc.ServiceRegistrar, srv ToggleCommandServiceServer) {
	s.RegisterService(&ToggleCommandService_ServiceDesc, srv)
}

func _ToggleCommandService_CreateToggle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToggleCommandServiceServer).CreateToggle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.indrasaputra.toggle.v1.ToggleCommandService/CreateToggle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToggleCommandServiceServer).CreateToggle(ctx, req.(*CreateToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToggleCommandService_EnableToggle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnableToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToggleCommandServiceServer).EnableToggle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.indrasaputra.toggle.v1.ToggleCommandService/EnableToggle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToggleCommandServiceServer).EnableToggle(ctx, req.(*EnableToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToggleCommandService_DisableToggle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisableToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToggleCommandServiceServer).DisableToggle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.indrasaputra.toggle.v1.ToggleCommandService/DisableToggle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToggleCommandServiceServer).DisableToggle(ctx, req.(*DisableToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToggleCommandService_DeleteToggle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToggleCommandServiceServer).DeleteToggle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.indrasaputra.toggle.v1.ToggleCommandService/DeleteToggle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToggleCommandServiceServer).DeleteToggle(ctx, req.(*DeleteToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ToggleCommandService_ServiceDesc is the grpc.ServiceDesc for ToggleCommandService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ToggleCommandService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.indrasaputra.toggle.v1.ToggleCommandService",
	HandlerType: (*ToggleCommandServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateToggle",
			Handler:    _ToggleCommandService_CreateToggle_Handler,
		},
		{
			MethodName: "EnableToggle",
			Handler:    _ToggleCommandService_EnableToggle_Handler,
		},
		{
			MethodName: "DisableToggle",
			Handler:    _ToggleCommandService_DisableToggle_Handler,
		},
		{
			MethodName: "DeleteToggle",
			Handler:    _ToggleCommandService_DeleteToggle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/indrasaputra/toggle/v1/toggle.proto",
}

// ToggleQueryServiceClient is the client API for ToggleQueryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ToggleQueryServiceClient interface {
	// Get a toggle.
	//
	// This endpoint gets a single toggle by its key.
	GetToggleByKey(ctx context.Context, in *GetToggleByKeyRequest, opts ...grpc.CallOption) (*GetToggleByKeyResponse, error)
	// Get many toggles.
	//
	// This endpoint gets all available toggles in the system.
	// Currently, it only retrieves 10 toggles at most.
	GetAllToggles(ctx context.Context, in *GetAllTogglesRequest, opts ...grpc.CallOption) (*GetAllTogglesResponse, error)
}

type toggleQueryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewToggleQueryServiceClient(cc grpc.ClientConnInterface) ToggleQueryServiceClient {
	return &toggleQueryServiceClient{cc}
}

func (c *toggleQueryServiceClient) GetToggleByKey(ctx context.Context, in *GetToggleByKeyRequest, opts ...grpc.CallOption) (*GetToggleByKeyResponse, error) {
	out := new(GetToggleByKeyResponse)
	err := c.cc.Invoke(ctx, "/proto.indrasaputra.toggle.v1.ToggleQueryService/GetToggleByKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *toggleQueryServiceClient) GetAllToggles(ctx context.Context, in *GetAllTogglesRequest, opts ...grpc.CallOption) (*GetAllTogglesResponse, error) {
	out := new(GetAllTogglesResponse)
	err := c.cc.Invoke(ctx, "/proto.indrasaputra.toggle.v1.ToggleQueryService/GetAllToggles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ToggleQueryServiceServer is the server API for ToggleQueryService service.
// All implementations must embed UnimplementedToggleQueryServiceServer
// for forward compatibility
type ToggleQueryServiceServer interface {
	// Get a toggle.
	//
	// This endpoint gets a single toggle by its key.
	GetToggleByKey(context.Context, *GetToggleByKeyRequest) (*GetToggleByKeyResponse, error)
	// Get many toggles.
	//
	// This endpoint gets all available toggles in the system.
	// Currently, it only retrieves 10 toggles at most.
	GetAllToggles(context.Context, *GetAllTogglesRequest) (*GetAllTogglesResponse, error)
	mustEmbedUnimplementedToggleQueryServiceServer()
}

// UnimplementedToggleQueryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedToggleQueryServiceServer struct {
}

func (UnimplementedToggleQueryServiceServer) GetToggleByKey(context.Context, *GetToggleByKeyRequest) (*GetToggleByKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetToggleByKey not implemented")
}
func (UnimplementedToggleQueryServiceServer) GetAllToggles(context.Context, *GetAllTogglesRequest) (*GetAllTogglesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllToggles not implemented")
}
func (UnimplementedToggleQueryServiceServer) mustEmbedUnimplementedToggleQueryServiceServer() {}

// UnsafeToggleQueryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ToggleQueryServiceServer will
// result in compilation errors.
type UnsafeToggleQueryServiceServer interface {
	mustEmbedUnimplementedToggleQueryServiceServer()
}

func RegisterToggleQueryServiceServer(s grpc.ServiceRegistrar, srv ToggleQueryServiceServer) {
	s.RegisterService(&ToggleQueryService_ServiceDesc, srv)
}

func _ToggleQueryService_GetToggleByKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetToggleByKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToggleQueryServiceServer).GetToggleByKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.indrasaputra.toggle.v1.ToggleQueryService/GetToggleByKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToggleQueryServiceServer).GetToggleByKey(ctx, req.(*GetToggleByKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ToggleQueryService_GetAllToggles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllTogglesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ToggleQueryServiceServer).GetAllToggles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.indrasaputra.toggle.v1.ToggleQueryService/GetAllToggles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ToggleQueryServiceServer).GetAllToggles(ctx, req.(*GetAllTogglesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ToggleQueryService_ServiceDesc is the grpc.ServiceDesc for ToggleQueryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ToggleQueryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.indrasaputra.toggle.v1.ToggleQueryService",
	HandlerType: (*ToggleQueryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetToggleByKey",
			Handler:    _ToggleQueryService_GetToggleByKey_Handler,
		},
		{
			MethodName: "GetAllToggles",
			Handler:    _ToggleQueryService_GetAllToggles_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/indrasaputra/toggle/v1/toggle.proto",
}
