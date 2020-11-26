package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/handlers"
	"github.com/lucthienbinh/golang_scem/middlewares"
)

var (
	router = gin.Default()
)

// RunServer is the entry point to start the server
func RunServer() {
	// gin.SetMode(gin.ReleaseMode)
	// export GIN_MODE=debug

	// Initial auth middleware
	middlewares.RunAuth()
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
	log.Print("Listening and serving HTTP on :5000")
	routeList()
	router.Run(":5000")
}

// RouteList will create our routes of our entire application
// this way every group of routes can be defined in their own file
// so this one won't be so messy
func routeList() {

	api := router.Group("/api")
	userRoutes(api)
	orderRoutes(api)

}
