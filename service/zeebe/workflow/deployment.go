package workflow

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

var (
	// ZbClient client to connect with zeebe engine
	zbClient zbc.Client
)

// ConnectZeebeEngine function
func ConnectZeebeEngine() {
	gatewayAddress := os.Getenv("BROKER_ADDRESS")
	newZbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         gatewayAddress,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	zbClient = newZbClient
}

// DeployNewWorkflow function
func DeployNewWorkflow() {

	ctx := context.Background()
	response, err := zbClient.NewDeployWorkflowCommand().AddResourceFile(os.Getenv("WORKFLOW_FILE_NAME_1")).Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.String())
}

// CreateNewInstance of workflow
func CreateNewInstance(orderID, customerReceiveID uint, payMethod string, useLongShip, useShortShip bool) uint {

	// After the workflow is deployed.
	variables := make(map[string]interface{})
	variables["order_id"] = orderID
	variables["customer_receive_id"] = customerReceiveID
	variables["pay_method"] = payMethod
	variables["use_long_ship"] = useLongShip
	variables["use_short_ship"] = useShortShip

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId(os.Getenv("WORKFLOW_ID_1")).LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	msg, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(msg.String())
	return uint(msg.WorkflowInstanceKey)
}
