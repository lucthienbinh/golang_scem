package workflow

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/lucthienbinh/golang_scem/internal/model"
	"github.com/lucthienbinh/golang_scem/internal/service/state_scem/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	address = "localhost:9001"
)

// DeployFullShipWorkflow function
func DeployFullShipWorkflow() error {

	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("STATE_SCEM_ADDRESS"), grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewStateScemServiceClient(conn)

	// Open our jsonFile
	jsonFile, err := os.Open(os.Getenv("FULL_SHIP_SS_FILE_1"))
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Source: https://tutorialedge.net/golang/parsing-json-with-golang/
	byteValue, _ := ioutil.ReadAll(jsonFile)
	message := &pb.DeployWorkflowRequest{}
	err = protojson.Unmarshal(byteValue, message)
	if err != nil {
		return err
	}

	r, err := c.DeployWorkflow(context.Background(), message)
	if err != nil {
		return err
	}
	log.Printf("Response from Server: %s", r)
	return nil

}

// DeployLongShipWorkflow function
func DeployLongShipWorkflow() error {

	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("STATE_SCEM_ADDRESS"), grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewStateScemServiceClient(conn)

	// Open our jsonFile
	jsonFile, err := os.Open(os.Getenv("LONG_SHIP_SS_FILE_1"))
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Source: https://tutorialedge.net/golang/parsing-json-with-golang/
	byteValue, _ := ioutil.ReadAll(jsonFile)
	message := &pb.DeployWorkflowRequest{}
	protojson.Unmarshal(byteValue, message)

	r, err := c.DeployWorkflow(context.Background(), message)
	if err != nil {
		return err
	}
	log.Printf("Response from Server: %s", r)
	return nil
}

// CreateFullShipInstance of workflow
func CreateFullShipInstance(orderWorkflowData *model.OrderWorkflowData) (string, uint, error) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("STATE_SCEM_ADDRESS"), grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		return "", uint(0), err
	}
	defer conn.Close()
	c := pb.NewStateScemServiceClient(conn)

	// After the workflow is deployed.
	variables := make(map[string]interface{})
	variables["order_id"] = orderWorkflowData.OrderID
	variables["shipper_receive_money"] = orderWorkflowData.ShipperReceiveMoney
	variables["use_long_ship"] = orderWorkflowData.UseLongShip
	variables["customer_receive_id"] = orderWorkflowData.CustomerReceiveID

	workflowVariables := make(map[string]interface{})
	workflowVariables["workflow_variable_list"] = variables
	workflowVariables["workflow_process_id"] = os.Getenv("FULL_SHIP_SS_ID_1")
	workflowVariablesJSON, _ := json.Marshal(workflowVariables)

	message := &pb.CreateWorkflowInstanceRequest{}
	protojson.Unmarshal(workflowVariablesJSON, message)

	r, err := c.CreateWorkflowInstance(context.Background(), message)
	if err != nil {
		return "", uint(0), err
	}
	log.Printf("Response from Server: %s", r)
	return r.WorkflowKey, uint(r.WorkflowInstanceID), nil
}
