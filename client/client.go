package main

import (
	"log"
	"time"

	pb "github.com/GoogleCloudPlatform/gcloud-golang-todos/todo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func main() {
	log.Printf("Before dial.")
	//	conn, err := grpc.Dial("localhost:5050")
	conn, err := grpc.Dial("192.168.59.103:5050")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	log.Printf("Dialed in.")
	defer conn.Close()
	c := pb.NewTodoServiceClient(conn)
	log.Printf("Made a client.")
	ctx, cancel := timeoutContext()
	defer cancel()
	t, err := c.NewTodo(ctx, &pb.NewTodoRequest{"foo"})
	if err != nil {
		log.Fatalf("Could not create new todo: %v", err)
	}
	log.Printf("Made todo: %v", t)
	ctx, cancel = timeoutContext()
	defer cancel()
	ts, err := c.ListTodos(ctx, &pb.NilRequest{})
	if err != nil {
		log.Fatalf("Could not list todos: %v", err)
	}
	log.Printf("Got todos: %+v", ts)
}
