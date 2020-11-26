package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/handlers"
)

func orderRoutes(rg *gin.RouterGroup) {

	transportType := rg.Group("/transport-type")
	// customers.Use(middlewares.ValidateSession)
	transportType.GET("/list", handlers.GetTransportTypeListHandler)
	transportType.GET("/id/:id", handlers.GetTransportTypeHandler)
	transportType.POST("/create", handlers.CreateTransportTypeHandler)
	transportType.PUT("/update/:id", handlers.UpdateTransportTypeHandler)
	transportType.DELETE("/delete/:id", handlers.DeleteTransportTypeHandler)

	order := rg.Group("/order")
	// employees.Use(middlewares.ValidateSession)
	order.GET("/list", handlers.GetOrderInfoListHandler)
	order.GET("/id/:id", handlers.GetOrderInfoHandler)
	order.POST("/create", handlers.CreateOrderInfoHandler)
	order.PUT("/update/:id", handlers.UpdateOrderInfoHandler)
	order.DELETE("/delete/:id", handlers.DeleteOrderInfoHandler)

}
