package zeebeclient

import (
	"context"
	"fmt"
	"log"

	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

// RunWorkerTaskA to start this worker
func RunWorkerTaskA() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         brokerAddr,
		UsePlaintextConnection: true,
	})
	if err != nil {
		panic(err)
	}

	go client.NewJobWorker().JobType("taskA").Handler(handleJob).Open()
}

func handleJob(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

	// headers, err := job.GetCustomHeadersAsMap()
	// if err != nil {
	// 	// failed to handle job as we require the custom job headers
	// 	failJob(client, job)
	// 	return
	// }

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	fmt.Println(variables)

	variables["workerStatus"] = 1
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)
	log.Println("Processing order:", variables["workerStatus"])
	// log.Println("Collect money using payment method:", headers["method"])

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
	if err != nil {
		panic(err)
	}
}
