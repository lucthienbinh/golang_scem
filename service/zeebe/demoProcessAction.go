package zeebeclient

import (
	"context"
	"fmt"

	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

const (
	brokerAddr       = "127.0.0.1:26500"
	workflowFileName = "demoProcess5.bpmn"
	workflowID       = "demoProcess5.bpmn"
)

// DeployNewWorkflow to zeebe
func DeployNewWorkflow() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         brokerAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	response, err := client.NewDeployWorkflowCommand().AddResourceFile(workflowFileName).Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.String())
}

// CreateNewInstance of workflow
func CreateNewInstance() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         brokerAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	// After the workflow is deployed.
	variables := make(map[string]interface{})
	variables["orderId"] = "31243"

	request, err := client.NewCreateInstanceCommand().BPMNProcessId(workflowID).LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	msg, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(msg.String())
}
