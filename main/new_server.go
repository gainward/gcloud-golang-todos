package main

import (
	"fmt"
	"io/ioutil"
	syslog "log"
	"net"
	"os"

	pb "github.com/GoogleCloudPlatform/gcloud-golang-todos/todo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
	"google.golang.org/grpc"
)

const (
	envProjID     = "GCLOUD_GOLANG_TODOS_PROJECT_ID"
	envPrivateKey = "GCLOUD_GOLANG_TODOS_KEY"
)

func getTodoGroupKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "TodoGroup", "", 1, nil)
}

func SaveTodo(c context.Context, t *pb.Todo) error {
	k := datastore.NewKey(c, "Todo", "", t.Id, getTodoGroupKey(c))
	k, err := datastore.Put(c, k, t)
	if err == nil {
		t.Id = k.ID()
	}
	return err
}

func Context(ctx context.Context, scopes ...string) context.Context {
	key, projID := os.Getenv(envPrivateKey), os.Getenv(envProjID)
	if key == "" || projID == "" {
		syslog.Fatalf("%v and %v must be set. See CONTRIBUTING.md",
			envProjID, envPrivateKey)
	}
	jsonKey, err := ioutil.ReadFile(key)
	if err != nil {
		syslog.Fatalf("Cannot read the JSON key file, err: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey, scopes...)
	if err != nil {
		syslog.Fatal(err)
	}
	_ = ctx
	return cloud.NewContext(projID, conf.Client(oauth2.NoContext))
	//	return cloud.WithContext(ctx, projID, conf.Client(oauth2.NoContext))
}

func datastoreContext(ctx context.Context) context.Context {
	return Context(ctx, datastore.ScopeDatastore, datastore.ScopeUserEmail)
}

type server struct{}

func (server) NewTodo(c context.Context, r *pb.NewTodoRequest) (*pb.Todo, error) {
	if r.Title == "" {
		return nil, fmt.Errorf("Title required.")
	}
	t := &pb.Todo{Title: r.Title}
	syslog.Printf("Saving new todo: %v", t)
	if err := SaveTodo(datastoreContext(c), t); err != nil {
		return nil, err
	}
	return t, nil
}

func (server) GetTodo(c context.Context, r *pb.GetTodoRequest) (*pb.Todo, error) {
	k := datastore.NewKey(c, "Todo", "", r.Id, getTodoGroupKey(c))
	todo := &pb.Todo{}
	err := datastore.Get(datastoreContext(c), k, todo)
	if err == nil {
		todo.Id = k.ID()
	}
	return todo, err
}

func (server) ListTodos(c context.Context, r *pb.NilRequest) (*pb.Todos, error) {
	var todos []*pb.Todo
	q := datastore.NewQuery("Todo").Ancestor(getTodoGroupKey(c))
	keys, err := q.GetAll(datastoreContext(c), &todos)
	if err == nil {
		for i := 0; i < len(keys); i++ {
			todos[i].Id = keys[i].ID()
		}
	}
	return &pb.Todos{Todos: todos}, err
}

func init() {
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		syslog.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &server{})
	go s.Serve(lis)
	syslog.Printf("Serving!")
}
