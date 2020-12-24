package worker

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

// RunBankPayment to start this worker
func RunOrderPayment() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         os.Getenv("BROKER_ADDRESS"),
		UsePlaintextConnection: true,
	})
	if err != nil {
		panic(err)
	}
	go client.NewJobWorker().JobType("order_payment").Handler(handleJob2).Open()
}

func handleJob2(client worker.JobClient, job entities.Job) {
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

	time.Sleep(10 * time.Second)
	_, ok := variables["order_id"].(uint)
	if ok == false {
		failJob(client, job)
	}
	payStatus, ok := variables["pay_status"].(bool)
	if ok == false {
		failJob(client, job)
	}
	payEmployeeID, ok := variables["pay_employee_id"].(uint)
	if ok == false {
		failJob(client, job)
	}
	var payServiceProvider = variables["pay_service_provider"]
	if (payStatus == false) || (payEmployeeID == 0 && payServiceProvider == nil) {
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
	log.Println("Processing order:", variables["workerStatus"])
	// log.Println("Collect money using payment method:", headers["method"])

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}
