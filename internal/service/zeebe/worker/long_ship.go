package worker

import (
	"context"
	"log"
	"os"

	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

// RunLongShip to start this worker
func RunLongShip() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         os.Getenv("BROKER_ADDRESS"),
		UsePlaintextConnection: true,
	})
	if err != nil {
		panic(err)
	}
	go client.NewJobWorker().JobType("long_ship").Handler(handleJobLongShip).Open()
}

func handleJobLongShip(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}
	// var uintOrderID uint
	// orderID, ok := variables["order_id"].(float64)
	// if ok == true {
	// 	uintOrderID = uint(orderID)
	// } else {
	// 	failJob(client, job)
	// 	return
	// }
	// var uintOrderShipID uint
	// orderShipID, ok := variables["order_ship_id"].(float64)
	// if ok == true {
	// 	uintOrderShipID = uint(orderShipID)
	// } else {
	// 	failJob(client, job)
	// 	return
	// }

	variables["long_ship_saved"] = true
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)
	log.Println("Processing order:", variables["order_id"])

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}
