package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lucthienbinh/golang_scem/api/middleware"
	"github.com/lucthienbinh/golang_scem/api/server"
	"github.com/lucthienbinh/golang_scem/internal/handler"
	ZBMessage "github.com/lucthienbinh/golang_scem/internal/service/zeebe/message"
	ZBWorker "github.com/lucthienbinh/golang_scem/internal/service/zeebe/worker"
	ZBWorkflow "github.com/lucthienbinh/golang_scem/internal/service/zeebe/workflow"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initial web auth middleware
	if os.Getenv("RUN_WEB_AUTH") == "yes" {
		middleware.RunWebAuth()
	}

	// Initial app auth middleware
	if os.Getenv("RUN_APP_AUTH") == "yes" {
		middleware.RunAppAuth()
	}

	// Select Postgres database
	if os.Getenv("SELECT_DATABASE") == "1" {
		if err := handler.ConnectPostgres(); err != nil {
			log.Print(err)
			os.Exit(1)
		}
		log.Print("Connected with posgres database!")
	}

	// Select MySQL database
	if os.Getenv("SELECT_DATABASE") == "2" {
		if err := handler.ConnectMySQL(); err != nil {
			log.Print(err)
			os.Exit(1)
		}
		log.Print("Connected with posgres database!")
	}

	if err := handler.RefreshDatabase(); err != nil {
		// if err := handlers.MigrationDatabase(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Print("Refreshed database!")

	ZBWorkflow.ConnectZeebeEngine()
	// ZBWorkflow.DeployNewWorkflow()
	// ZBWorkflow.CreateNewInstance(1, 2, "cash", true, true)

	ZBMessage.ConnectZeebeEngine()
	// ZBMessage.MoneyReceived(1, 1)

	ZBWorker.RunBankPayment()
	ZBWorker.RunSavePayment()

	// Our servers will live in the routes package
	server.RunServer()
}
