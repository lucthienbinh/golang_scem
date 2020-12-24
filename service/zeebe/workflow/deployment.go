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
	ZbClient zbc.Client
)

func connectZeebeEngine() {
	gatewayAddress := os.Getenv("BROKER_ADDRESS")
	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         gatewayAddress,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	ZbClient = zbClient
}

func deployNewWorkflow() {

	ctx := context.Background()
	response, err := ZbClient.NewDeployWorkflowCommand().AddResourceFile(os.Getenv("WORKFLOW_FILE_NAME_1")).Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.String())
}

// RunZeebeService to coonect zeebe and deploy workflow
func RunZeebeService() {
	connectZeebeEngine()
	deployNewWorkflow()
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

	request, err := ZbClient.NewCreateInstanceCommand().BPMNProcessId(os.Getenv("WORKFLOW_ID_1")).LatestVersion().VariablesFromMap(variables)
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
