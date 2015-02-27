// Code generated by protoc-gen-go.
// source: todo.proto
// DO NOT EDIT!

/*
Package todo is a generated protocol buffer package.

It is generated from these files:
	todo.proto

It has these top-level messages:
	NilRequest
	NewTodoRequest
	GetTodoRequest
	Todo
	Todos
*/
package todo

import proto "github.com/golang/protobuf/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type NilRequest struct {
}

func (m *NilRequest) Reset()         { *m = NilRequest{} }
func (m *NilRequest) String() string { return proto.CompactTextString(m) }
func (*NilRequest) ProtoMessage()    {}

type NewTodoRequest struct {
	Title string `protobuf:"bytes,1,opt,name=title" json:"title,omitempty"`
}

func (m *NewTodoRequest) Reset()         { *m = NewTodoRequest{} }
func (m *NewTodoRequest) String() string { return proto.CompactTextString(m) }
func (*NewTodoRequest) ProtoMessage()    {}

type GetTodoRequest struct {
	Id int64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *GetTodoRequest) Reset()         { *m = GetTodoRequest{} }
func (m *GetTodoRequest) String() string { return proto.CompactTextString(m) }
func (*GetTodoRequest) ProtoMessage()    {}

type Todo struct {
	Title     string `protobuf:"bytes,1,opt,name=title" json:"title,omitempty"`
	Id        int64  `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
	Completed bool   `protobuf:"varint,3,opt,name=completed" json:"completed,omitempty"`
}

func (m *Todo) Reset()         { *m = Todo{} }
func (m *Todo) String() string { return proto.CompactTextString(m) }
func (*Todo) ProtoMessage()    {}

type Todos struct {
	Todos []*Todo `protobuf:"bytes,1,rep,name=todos" json:"todos,omitempty"`
}

func (m *Todos) Reset()         { *m = Todos{} }
func (m *Todos) String() string { return proto.CompactTextString(m) }
func (*Todos) ProtoMessage()    {}

func (m *Todos) GetTodos() []*Todo {
	if m != nil {
		return m.Todos
	}
	return nil
}

func init() {
}

// Client API for TodoService service

type TodoServiceClient interface {
	NewTodo(ctx context.Context, in *NewTodoRequest, opts ...grpc.CallOption) (*Todo, error)
	GetTodo(ctx context.Context, in *GetTodoRequest, opts ...grpc.CallOption) (*Todo, error)
	ListTodos(ctx context.Context, in *NilRequest, opts ...grpc.CallOption) (*Todos, error)
}

type todoServiceClient struct {
	cc *grpc.ClientConn
}

func NewTodoServiceClient(cc *grpc.ClientConn) TodoServiceClient {
	return &todoServiceClient{cc}
}

func (c *todoServiceClient) NewTodo(ctx context.Context, in *NewTodoRequest, opts ...grpc.CallOption) (*Todo, error) {
	out := new(Todo)
	err := grpc.Invoke(ctx, "/todo.TodoService/NewTodo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) GetTodo(ctx context.Context, in *GetTodoRequest, opts ...grpc.CallOption) (*Todo, error) {
	out := new(Todo)
	err := grpc.Invoke(ctx, "/todo.TodoService/GetTodo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) ListTodos(ctx context.Context, in *NilRequest, opts ...grpc.CallOption) (*Todos, error) {
	out := new(Todos)
	err := grpc.Invoke(ctx, "/todo.TodoService/ListTodos", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TodoService service

type TodoServiceServer interface {
	NewTodo(context.Context, *NewTodoRequest) (*Todo, error)
	GetTodo(context.Context, *GetTodoRequest) (*Todo, error)
	ListTodos(context.Context, *NilRequest) (*Todos, error)
}

func RegisterTodoServiceServer(s *grpc.Server, srv TodoServiceServer) {
	s.RegisterService(&_TodoService_serviceDesc, srv)
}

func _TodoService_NewTodo_Handler(srv interface{}, ctx context.Context, buf []byte) (proto.Message, error) {
	in := new(NewTodoRequest)
	if err := proto.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(TodoServiceServer).NewTodo(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _TodoService_GetTodo_Handler(srv interface{}, ctx context.Context, buf []byte) (proto.Message, error) {
	in := new(GetTodoRequest)
	if err := proto.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(TodoServiceServer).GetTodo(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _TodoService_ListTodos_Handler(srv interface{}, ctx context.Context, buf []byte) (proto.Message, error) {
	in := new(NilRequest)
	if err := proto.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(TodoServiceServer).ListTodos(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _TodoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "todo.TodoService",
	HandlerType: (*TodoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewTodo",
			Handler:    _TodoService_NewTodo_Handler,
		},
		{
			MethodName: "GetTodo",
			Handler:    _TodoService_GetTodo_Handler,
		},
		{
			MethodName: "ListTodos",
			Handler:    _TodoService_ListTodos_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
