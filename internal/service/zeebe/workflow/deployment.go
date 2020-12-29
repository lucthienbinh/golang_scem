package workflow

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/lucthienbinh/golang_scem/internal/model"
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
func CreateNewInstance(orderWorkflowData *model.OrderWorkflowData) (uint, uint, error) {

	// After the workflow is deployed.
	variables := make(map[string]interface{})
	variables["order_id"] = orderWorkflowData.OrderID
	variables["order_pay_id"] = orderWorkflowData.OrderPayID
	variables["pay_method"] = orderWorkflowData.PayMethod
	variables["shipper_receive_money"] = orderWorkflowData.ShipperReceiveMoney
	variables["use_long_ship"] = orderWorkflowData.UseLongShip
	variables["use_short_ship"] = orderWorkflowData.UseShortShip
	variables["customer_receive_id"] = orderWorkflowData.CustomerReceiveID

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId(os.Getenv("WORKFLOW_ID_1")).LatestVersion().VariablesFromMap(variables)
	if err != nil {
		return uint(0), uint(0), err
	}

	ctx := context.Background()
	msg, err := request.Send(ctx)
	if err != nil {
		return uint(0), uint(0), err
	}
	log.Println(msg.String())
	return uint(msg.WorkflowKey), uint(msg.WorkflowInstanceKey), nil
}
