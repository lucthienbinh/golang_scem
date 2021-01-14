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
	if os.Getenv("RUNENV") != "docker" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Initial web auth middleware
	if os.Getenv("RUN_WEB_AUTH") == "yes" {
		runWebAuth()
	}

	// Initial app auth middleware
	if os.Getenv("RUN_APP_AUTH") == "yes" {
		runAppAuth()
	}

	// Select database
	if os.Getenv("SELECT_DATABASE") == "1" {
		connectPostgress()
	} else if os.Getenv("SELECT_DATABASE") == "2" {
		connectMySQL()
	} else if os.Getenv("SELECT_DATABASE") == "3" {
		connectSQLite()
	} else {
		log.Println("No database selected!")
		os.Exit(1)
	}

	if err := handler.RefreshDatabase(); err != nil {
		// if err := handler.MigrationDatabase(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Print("Database refreshed!")

	if os.Getenv("STATE_SERVICE") == "1" {
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

func connectSQLite() {
	if err := handler.ConnectSQLite(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("Connected with sqlite database!")
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
	ZBWorker.RunLongShip()
	ZBWorker.RunShortShip()
	ZBWorker.RunLongShipFinish()
}
