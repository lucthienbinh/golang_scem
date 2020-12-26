package main

import (
	"context"
	"log"

	"github.com/appleboy/gorush/rpc/proto"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc"
)

const (
	address = "localhost:9000"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGorushClient(conn)

	r, err := c.Send(context.Background(), &proto.NotificationRequest{
		Platform: 2,
		Tokens:   []string{"f09fih-jQ9GhMz3riGfmJv:APA91bH7opzi3nvIeY1GLvJb0zZClx19ZztB5-6Bgg4jIsBi-9fnZWHpqYo1Za78W93VbdyiQureIFkck0MA6AaFik7LwQ2gIburmRCV2eR4ZBIp-YjQKRhIUHAYbu6YyQmfEJPsDgbn"},
		Message:  "Your package will arrive soon!",
		Title:    "Attention!",
		Sound:    "test",
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"channelId": {
					Kind: &structpb.Value_StringValue{StringValue: "channel_id_1"},
				},
				"key2": {
					Kind: &structpb.Value_NumberValue{NumberValue: 2},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Success: %t\n", r.Success)
	log.Printf("Count: %d\n", r.Counts)
}
