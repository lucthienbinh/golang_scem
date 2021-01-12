package workflow

import (
	"context"
	"log"
	"time"

	"github.com/lucthienbinh/golang_scem/internal/service/state_scem/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

var (
	stateScemClient pb.StateScemServiceClient
)

// ConnectGolangStateScem function
func ConnectGolangStateScem() error {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("OK")
	stateScemClient = pb.NewStateScemServiceClient(conn)

	// keep alive in 3 mins
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	// defer cancel()
	return nil
}

const (
	address = "localhost:9001"
)

// DeployFullShipWorkflow function
func DeployFullShipWorkflow() error {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewStateScemServiceClient(conn)

	workflowModelPB := []*pb.WorkflowModel{{
		WorkflowProcessID: "full_ship_process_1",
		Step:              1,
		Type:              1,
		Name:              "Start",
		NextStep1:         2,
	}, {}}
	message := pb.DeployWorkflowRequest{WorkflowModel: workflowModelPB}
	r, err := c.DeployWorkflow(context.Background(), &message)
	if err != nil {
		return err
	}
	log.Printf("Response from Server: %s", r)
	return nil

}

// DeployLongShipWorkflow function
func DeployLongShipWorkflow() {
	message := pb.DeployWorkflowRequest{}
	response, err := stateScemClient.DeployWorkflow(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from Server: %s", response)
}
