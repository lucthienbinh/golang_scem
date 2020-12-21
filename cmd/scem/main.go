package main

import (
	"log"
	"os"

	"github.com/lucthienbinh/golang_scem/api/server"
	"github.com/lucthienbinh/golang_scem/internal/handler"
)

func main() {
	// Initial web auth middleware
	// middleware.RunWebAuth()

	// Initial app auth middleware
	// middleware.RunAppAuth()

	// Connect Postgres database
	// if err := handler.ConnectPostgres(); err != nil {
	if err := handler.ConnectMySQL(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Print("Connected with posgres database with: user=postgres dbname=scem_database port=543")
	if err := handler.RefreshDatabase(); err != nil {
		// if err := handlers.MigrationDatabase(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Print("Refreshed database!")
	// Our server will live in the routes package
	server.RunServer()
}
