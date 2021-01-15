package worker

import (
	"context"
	"log"
	"os"

	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"

	CommonService "github.com/lucthienbinh/golang_scem/internal/service/common"
)

// RunOrderShortShip to start this worker
func RunOrderShortShip() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         os.Getenv("BROKER_ADDRESS"),
		UsePlaintextConnection: true,
	})
	if err != nil {
		panic(err)
	}
	go client.NewJobWorker().JobType("order_short_ship").Handler(handleJobOrderShortShip).Open()
}

func handleJobOrderShortShip(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}
	var uintOrderID uint
	orderID, ok := variables["order_id"].(float64)
	if ok == true {
		uintOrderID = uint(orderID)
	} else {
		failJob(client, job)
		return
	}
	orderShortShipID, err := CommonService.CreateOrderShortShip(uintOrderID)
	if err != nil {
		failJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)
	log.Println("Created short ship id:", orderShortShipID)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}
