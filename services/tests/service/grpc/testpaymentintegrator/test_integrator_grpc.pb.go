// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: test_integrator.proto

package testpaymentintegrator

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

// TestPaymentIntegratorClient is the client API for TestPaymentIntegrator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestPaymentIntegratorClient interface {
	PublishTest(ctx context.Context, in *PublishAction, opts ...grpc.CallOption) (*Test, error)
	QualifyUserTestAndSession(ctx context.Context, in *SessionAction, opts ...grpc.CallOption) (*TestSession, error)
	InitiateTestSession(ctx context.Context, in *SessionAction, opts ...grpc.CallOption) (*TestSession, error)
}

type testPaymentIntegratorClient struct {
	cc grpc.ClientConnInterface
}

func NewTestPaymentIntegratorClient(cc grpc.ClientConnInterface) TestPaymentIntegratorClient {
	return &testPaymentIntegratorClient{cc}
}

func (c *testPaymentIntegratorClient) PublishTest(ctx context.Context, in *PublishAction, opts ...grpc.CallOption) (*Test, error) {
	out := new(Test)
	err := c.cc.Invoke(ctx, "/testpaymentintegrator.TestPaymentIntegrator/PublishTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testPaymentIntegratorClient) QualifyUserTestAndSession(ctx context.Context, in *SessionAction, opts ...grpc.CallOption) (*TestSession, error) {
	out := new(TestSession)
	err := c.cc.Invoke(ctx, "/testpaymentintegrator.TestPaymentIntegrator/QualifyUserTestAndSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testPaymentIntegratorClient) InitiateTestSession(ctx context.Context, in *SessionAction, opts ...grpc.CallOption) (*TestSession, error) {
	out := new(TestSession)
	err := c.cc.Invoke(ctx, "/testpaymentintegrator.TestPaymentIntegrator/InitiateTestSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TestPaymentIntegratorServer is the server API for TestPaymentIntegrator service.
// All implementations must embed UnimplementedTestPaymentIntegratorServer
// for forward compatibility
type TestPaymentIntegratorServer interface {
	PublishTest(context.Context, *PublishAction) (*Test, error)
	QualifyUserTestAndSession(context.Context, *SessionAction) (*TestSession, error)
	InitiateTestSession(context.Context, *SessionAction) (*TestSession, error)
	mustEmbedUnimplementedTestPaymentIntegratorServer()
}

// UnimplementedTestPaymentIntegratorServer must be embedded to have forward compatible implementations.
type UnimplementedTestPaymentIntegratorServer struct {
}

func (UnimplementedTestPaymentIntegratorServer) PublishTest(context.Context, *PublishAction) (*Test, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishTest not implemented")
}
func (UnimplementedTestPaymentIntegratorServer) QualifyUserTestAndSession(context.Context, *SessionAction) (*TestSession, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QualifyUserTestAndSession not implemented")
}
func (UnimplementedTestPaymentIntegratorServer) InitiateTestSession(context.Context, *SessionAction) (*TestSession, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitiateTestSession not implemented")
}
func (UnimplementedTestPaymentIntegratorServer) mustEmbedUnimplementedTestPaymentIntegratorServer() {}

// UnsafeTestPaymentIntegratorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestPaymentIntegratorServer will
// result in compilation errors.
type UnsafeTestPaymentIntegratorServer interface {
	mustEmbedUnimplementedTestPaymentIntegratorServer()
}

func RegisterTestPaymentIntegratorServer(s grpc.ServiceRegistrar, srv TestPaymentIntegratorServer) {
	s.RegisterService(&TestPaymentIntegrator_ServiceDesc, srv)
}

func _TestPaymentIntegrator_PublishTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishAction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestPaymentIntegratorServer).PublishTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/testpaymentintegrator.TestPaymentIntegrator/PublishTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestPaymentIntegratorServer).PublishTest(ctx, req.(*PublishAction))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestPaymentIntegrator_QualifyUserTestAndSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionAction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestPaymentIntegratorServer).QualifyUserTestAndSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/testpaymentintegrator.TestPaymentIntegrator/QualifyUserTestAndSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestPaymentIntegratorServer).QualifyUserTestAndSession(ctx, req.(*SessionAction))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestPaymentIntegrator_InitiateTestSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionAction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestPaymentIntegratorServer).InitiateTestSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/testpaymentintegrator.TestPaymentIntegrator/InitiateTestSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestPaymentIntegratorServer).InitiateTestSession(ctx, req.(*SessionAction))
	}
	return interceptor(ctx, in, info, handler)
}

// TestPaymentIntegrator_ServiceDesc is the grpc.ServiceDesc for TestPaymentIntegrator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TestPaymentIntegrator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "testpaymentintegrator.TestPaymentIntegrator",
	HandlerType: (*TestPaymentIntegratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishTest",
			Handler:    _TestPaymentIntegrator_PublishTest_Handler,
		},
		{
			MethodName: "QualifyUserTestAndSession",
			Handler:    _TestPaymentIntegrator_QualifyUserTestAndSession_Handler,
		},
		{
			MethodName: "InitiateTestSession",
			Handler:    _TestPaymentIntegrator_InitiateTestSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test_integrator.proto",
}
