package main

import (
	"log"
	"os"

	"github.com/lucthienbinh/golang_scem/handlers"
	"github.com/lucthienbinh/golang_scem/middlewares"
	"github.com/lucthienbinh/golang_scem/routes"
)

func main() {
	// Initial web auth middleware
	middlewares.RunWebAuth()

	// Initial app auth middleware
	middlewares.RunAppAuth()

	// Connect Postgres database
	if err := handlers.ConnectPostgres(); err != nil {
		// if err := handlers.ConnectMySQL(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	if err := handlers.RefreshDatabase(); err != nil {
		// if err := handlers.MigrationDatabase(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Print("Connected with posgres database with: user=postgres dbname=scem_database port=543")
	// Our server will live in the routes package
	routes.RunServer()
}
