package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/handler"
)

// OrderRoutes to manage order model
func WebOrderRoutes(rg *gin.RouterGroup) {

	transportType := rg.Group("/transport-type")
	transportType.GET("/list", handler.GetTransportTypeListHandler)
	transportType.GET("/id/:id", handler.GetTransportTypeHandler)
	transportType.POST("/create", handler.CreateTransportTypeHandler)
	transportType.PUT("/update/:id", handler.UpdateTransportTypeHandler)
	transportType.DELETE("/delete/:id", handler.DeleteTransportTypeHandler)

	order := rg.Group("/order")
	order.GET("/list", handler.GetOrderInfoListHandler)
	order.GET("/id/:id", handler.GetOrderInfoHandler)
	order.POST("/create", handler.CreateOrderInfoHandler)
	order.PUT("/update/:id", handler.UpdateOrderInfoHandler)
	order.DELETE("/delete/:id", handler.DeleteOrderInfoHandler)

}
