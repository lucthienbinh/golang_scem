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
		runWebAuth()
	}

	// Initial app auth middleware
	if os.Getenv("RUN_APP_AUTH") == "yes" {
		runAppAuth()
	}

	// Select Postgres database
	if os.Getenv("SELECT_DATABASE") == "1" {
		connectPostgress()
	}

	// Select MySQL database
	if os.Getenv("SELECT_DATABASE") == "2" {
		connectMySQL()
	}

	if err := handler.RefreshDatabase(); err != nil {
		// if err := handlers.MigrationDatabase(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Print("Refreshed database!")

	if os.Getenv("USE_ZEEBE") == "1" {
		connectZeebeClient()
	}

	// Our servers will live in the routes package
	server.RunServer()
}

// Source code: https://www.devdungeon.com/content/working-files-go#read_all
func runWebAuth() {
	sessionKey := []byte(os.Getenv("SESSION_KEY"))
	if err := middleware.RunWebAuth(sessionKey); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Web authenticate activated!")
}

func runAppAuth() {
	if err := middleware.RunAppAuth(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("App authenticate activated!")
}

func connectPostgress() {
	if err := handler.ConnectPostgres(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Connected with posgres database!")
}

func connectMySQL() {
	if err := handler.ConnectMySQL(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("Connected with posgres database!")
}

func connectZeebeClient() {
	if err := ZBWorkflow.ConnectZeebeEngine(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("Zeebe workflow package connected with zeebe!")
	if err := ZBMessage.ConnectZeebeEngine(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("Zeebe message package connected with zeebe!")
	// Run Zebee service
	ZBWorker.RunCreditPayment()
	ZBWorker.RunLongShip()
	ZBWorker.RunShortShip()
	ZBWorker.RunLongShipFinish()
}
