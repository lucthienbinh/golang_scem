package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/handlers"
	"github.com/lucthienbinh/golang_scem/middlewares"
)

func webAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/web/loginJSON", handlers.WebLoginHandler)
	rg.GET("/web/logout", handlers.WebLogoutHandler)
}

func appAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/app/loginJSON", handlers.AppLoginHandler)
	rg.GET("/app/logout", handlers.AppLogoutHandler)
	// validate old token
	accessToken := rg.Group("/app/access-token", middlewares.ValidateAppTokenForRefresh())
	accessToken.POST("/get-new", handlers.AppReloginHandler)
	accessToken.GET("/check-old", handlers.AppOpenHandler)
}

func userRoutes(rg *gin.RouterGroup) {

	customer := rg.Group("/customer")
	customer.GET("/list", handlers.GetCustomerListHandler)
	customer.GET("/id/:id", handlers.GetCustomerHandler)
	customer.POST("/create", handlers.CreateCustomerHandler)
	customer.PUT("/update/:id", handlers.UpdateCustomerHandler)
	customer.DELETE("/delete/:id", handlers.DeleteCustomerHandler)

	employee := rg.Group("/employee")
	employee.GET("/list", handlers.GetEmployeeListHandler)
	employee.GET("/id/:id", handlers.GetEmployeeHandler)
	employee.POST("/create", handlers.CreateEmployeeHandler)
	employee.POST("/upload/image", handlers.ImageEmployeeHandler)
	employee.PUT("/update/:id", handlers.UpdateEmployeeHandler)
	employee.DELETE("/delete/:id", handlers.DeleteEmployeeHandler)

	employeeType := rg.Group("/employee-type")
	employeeType.GET("/list", handlers.GetEmployeeTypeListHandler)
	employeeType.GET("/id/:id", handlers.GetEmployeeTypeHandler)
	employeeType.POST("/create", handlers.CreateEmployeeTypeHandler)
	employeeType.PUT("/update/:id", handlers.UpdateEmployeeTypeHandler)
	employeeType.DELETE("/delete/:id", handlers.DeleteEmployeeTypeHandler)

	// middlewares.ValidateAppToken()
	deliveryLocation := rg.Group("/delivery-location")
	deliveryLocation.GET("/list", handlers.GetDeliveryLocationListHandler)
	deliveryLocation.GET("/id/:id", handlers.GetDeliveryLocationHandler)
	deliveryLocation.POST("/create", handlers.CreateDeliveryLocationHandler)
	deliveryLocation.PUT("/update/:id", handlers.UpdateDeliveryLocationHandler)
	deliveryLocation.DELETE("/delete/:id", handlers.DeleteDeliveryLocationHandler)

}
