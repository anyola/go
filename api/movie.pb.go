package api

import (
	context "context"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = proto.ProtoPackageIsVersion4

type Movie struct {
	ID       int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title    string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Year     int32  `protobuf:"varint,3,opt,name=year,proto3" json:"year,omitempty"`
	Director string `protobuf:"bytes,4,opt,name=director,proto3" json:"director,omitempty"`
}
type LoadMoviesRequest struct {
	From int32 `protobuf:"varint,1,opt,name=from,proto3" json:"from,omitempty"`
	To   int32 `protobuf:"varint,2,opt,name=to,proto3" json:"to,omitempty"`
}
type LoadError struct {
	ID      int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}
type LoadMoviesResponse struct {
	Movies []*Movie     `protobuf:"bytes,1,rep,name=movies,proto3" json:"movies,omitempty"`
	Errors []*LoadError `protobuf:"bytes,2,rep,name=errors,proto3" json:"errors,omitempty"`
}
type GetLastMovieRequest struct{}

func (*Movie) Reset()                 {}
func (*Movie) String() string         { return "" }
func (*Movie) ProtoMessage()          {}
func (*LoadMoviesRequest) Reset()     {}
func (*LoadMoviesRequest) String() string { return "" }
func (*LoadMoviesRequest) ProtoMessage() {}
func (*LoadError) Reset()             {}
func (*LoadError) String() string     { return "" }
func (*LoadError) ProtoMessage()      {}
func (*LoadMoviesResponse) Reset()    {}
func (*LoadMoviesResponse) String() string { return "" }
func (*LoadMoviesResponse) ProtoMessage() {}
func (*GetLastMovieRequest) Reset()   {}
func (*GetLastMovieRequest) String() string { return "" }
func (*GetLastMovieRequest) ProtoMessage() {}

type MovieServiceClient interface {
	LoadMovies(context.Context, *LoadMoviesRequest, ...grpc.CallOption) (*LoadMoviesResponse, error)
	GetLastMovie(context.Context, *GetLastMovieRequest, ...grpc.CallOption) (*Movie, error)
}
type movieServiceClient struct{ cc *grpc.ClientConn }

func NewMovieServiceClient(cc *grpc.ClientConn) MovieServiceClient {
	return &movieServiceClient{cc}
}
func (c *movieServiceClient) LoadMovies(ctx context.Context, in *LoadMoviesRequest, opts ...grpc.CallOption) (*LoadMoviesResponse, error) {
	out := new(LoadMoviesResponse)
	err := c.cc.Invoke(ctx, "/movieapi.MovieService/LoadMovies", in, out, opts...)
	return out, err
}
func (c *movieServiceClient) GetLastMovie(ctx context.Context, in *GetLastMovieRequest, opts ...grpc.CallOption) (*Movie, error) {
	out := new(Movie)
	err := c.cc.Invoke(ctx, "/movieapi.MovieService/GetLastMovie", in, out, opts...)
	return out, err
}

type MovieServiceServer interface {
	LoadMovies(context.Context, *LoadMoviesRequest) (*LoadMoviesResponse, error)
	GetLastMovie(context.Context, *GetLastMovieRequest) (*Movie, error)
}
type UnimplementedMovieServiceServer struct{}
func (UnimplementedMovieServiceServer) LoadMovies(context.Context, *LoadMoviesRequest) (*LoadMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadMovies not implemented")
}
func (UnimplementedMovieServiceServer) GetLastMovie(context.Context, *GetLastMovieRequest) (*Movie, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastMovie not implemented")
}
func RegisterMovieServiceServer(s *grpc.Server, srv MovieServiceServer) {
	s.RegisterService(&MovieService_ServiceDesc, srv)
}
func _MovieService_LoadMovies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadMoviesRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MovieServiceServer).LoadMovies(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/movieapi.MovieService/LoadMovies"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MovieServiceServer).LoadMovies(ctx, req.(*LoadMoviesRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _MovieService_GetLastMovie_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLastMovieRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(MovieServiceServer).GetLastMovie(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/movieapi.MovieService/GetLastMovie"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(MovieServiceServer).GetLastMovie(ctx, req.(*GetLastMovieRequest)) }
	return interceptor(ctx, in, info, handler)
}
var MovieService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "movieapi.MovieService",
	HandlerType: (*MovieServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "LoadMovies", Handler: _MovieService_LoadMovies_Handler},
		{MethodName: "GetLastMovie", Handler: _MovieService_GetLastMovie_Handler},
	},
	Streams: []grpc.StreamDesc{},
	Metadata: "api/movie.proto",
}
