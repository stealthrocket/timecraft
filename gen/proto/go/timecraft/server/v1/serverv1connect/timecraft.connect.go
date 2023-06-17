// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: timecraft/server/v1/timecraft.proto

package serverv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/stealthrocket/timecraft/gen/proto/go/timecraft/server/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// TimecraftServiceName is the fully-qualified name of the TimecraftService service.
	TimecraftServiceName = "timecraft.server.v1.TimecraftService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// TimecraftServiceParentProcedure is the fully-qualified name of the TimecraftService's Parent RPC.
	TimecraftServiceParentProcedure = "/timecraft.server.v1.TimecraftService/Parent"
	// TimecraftServiceSpawnProcedure is the fully-qualified name of the TimecraftService's Spawn RPC.
	TimecraftServiceSpawnProcedure = "/timecraft.server.v1.TimecraftService/Spawn"
	// TimecraftServiceVersionProcedure is the fully-qualified name of the TimecraftService's Version
	// RPC.
	TimecraftServiceVersionProcedure = "/timecraft.server.v1.TimecraftService/Version"
)

// TimecraftServiceClient is a client for the timecraft.server.v1.TimecraftService service.
type TimecraftServiceClient interface {
	Parent(context.Context, *connect_go.Request[v1.ParentRequest]) (*connect_go.Response[v1.ParentResponse], error)
	Spawn(context.Context, *connect_go.Request[v1.SpawnRequest]) (*connect_go.Response[v1.SpawnResponse], error)
	Version(context.Context, *connect_go.Request[v1.VersionRequest]) (*connect_go.Response[v1.VersionResponse], error)
}

// NewTimecraftServiceClient constructs a client for the timecraft.server.v1.TimecraftService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTimecraftServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) TimecraftServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &timecraftServiceClient{
		parent: connect_go.NewClient[v1.ParentRequest, v1.ParentResponse](
			httpClient,
			baseURL+TimecraftServiceParentProcedure,
			opts...,
		),
		spawn: connect_go.NewClient[v1.SpawnRequest, v1.SpawnResponse](
			httpClient,
			baseURL+TimecraftServiceSpawnProcedure,
			opts...,
		),
		version: connect_go.NewClient[v1.VersionRequest, v1.VersionResponse](
			httpClient,
			baseURL+TimecraftServiceVersionProcedure,
			opts...,
		),
	}
}

// timecraftServiceClient implements TimecraftServiceClient.
type timecraftServiceClient struct {
	parent  *connect_go.Client[v1.ParentRequest, v1.ParentResponse]
	spawn   *connect_go.Client[v1.SpawnRequest, v1.SpawnResponse]
	version *connect_go.Client[v1.VersionRequest, v1.VersionResponse]
}

// Parent calls timecraft.server.v1.TimecraftService.Parent.
func (c *timecraftServiceClient) Parent(ctx context.Context, req *connect_go.Request[v1.ParentRequest]) (*connect_go.Response[v1.ParentResponse], error) {
	return c.parent.CallUnary(ctx, req)
}

// Spawn calls timecraft.server.v1.TimecraftService.Spawn.
func (c *timecraftServiceClient) Spawn(ctx context.Context, req *connect_go.Request[v1.SpawnRequest]) (*connect_go.Response[v1.SpawnResponse], error) {
	return c.spawn.CallUnary(ctx, req)
}

// Version calls timecraft.server.v1.TimecraftService.Version.
func (c *timecraftServiceClient) Version(ctx context.Context, req *connect_go.Request[v1.VersionRequest]) (*connect_go.Response[v1.VersionResponse], error) {
	return c.version.CallUnary(ctx, req)
}

// TimecraftServiceHandler is an implementation of the timecraft.server.v1.TimecraftService service.
type TimecraftServiceHandler interface {
	Parent(context.Context, *connect_go.Request[v1.ParentRequest]) (*connect_go.Response[v1.ParentResponse], error)
	Spawn(context.Context, *connect_go.Request[v1.SpawnRequest]) (*connect_go.Response[v1.SpawnResponse], error)
	Version(context.Context, *connect_go.Request[v1.VersionRequest]) (*connect_go.Response[v1.VersionResponse], error)
}

// NewTimecraftServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTimecraftServiceHandler(svc TimecraftServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(TimecraftServiceParentProcedure, connect_go.NewUnaryHandler(
		TimecraftServiceParentProcedure,
		svc.Parent,
		opts...,
	))
	mux.Handle(TimecraftServiceSpawnProcedure, connect_go.NewUnaryHandler(
		TimecraftServiceSpawnProcedure,
		svc.Spawn,
		opts...,
	))
	mux.Handle(TimecraftServiceVersionProcedure, connect_go.NewUnaryHandler(
		TimecraftServiceVersionProcedure,
		svc.Version,
		opts...,
	))
	return "/timecraft.server.v1.TimecraftService/", mux
}

// UnimplementedTimecraftServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTimecraftServiceHandler struct{}

func (UnimplementedTimecraftServiceHandler) Parent(context.Context, *connect_go.Request[v1.ParentRequest]) (*connect_go.Response[v1.ParentResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("timecraft.server.v1.TimecraftService.Parent is not implemented"))
}

func (UnimplementedTimecraftServiceHandler) Spawn(context.Context, *connect_go.Request[v1.SpawnRequest]) (*connect_go.Response[v1.SpawnResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("timecraft.server.v1.TimecraftService.Spawn is not implemented"))
}

func (UnimplementedTimecraftServiceHandler) Version(context.Context, *connect_go.Request[v1.VersionRequest]) (*connect_go.Response[v1.VersionResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("timecraft.server.v1.TimecraftService.Version is not implemented"))
}
