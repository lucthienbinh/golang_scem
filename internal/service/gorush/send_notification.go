package gorush

import (
	"context"
	"log"
	"os"

	"github.com/appleboy/gorush/rpc/proto"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc"
)

// Client to connect gorush
func Client(tokenString, title, message string) error {
	address := os.Getenv("GORUSH_ADDRESS")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	c := proto.NewGorushClient(conn)
	log.Println("token:", tokenString)
	r, err := c.Send(context.Background(), &proto.NotificationRequest{
		Platform: 2,
		Tokens:   []string{tokenString},
		Message:  message,
		Title:    title,
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
		return err
	}
	log.Printf("Success: %t\n", r.Success)
	log.Printf("Count: %d\n", r.Counts)
	return nil
}
