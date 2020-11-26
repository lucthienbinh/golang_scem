package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/handlers"
)

func userRoutes(rg *gin.RouterGroup) {
	// Authorization user
	authUsers := rg.Group("/auth-user")
	authUsers.POST("/loginJSON", handlers.LoginHandler)
	authUsers.GET("/logout", handlers.LogoutHandler)

	customer := rg.Group("/customer")
	// customers.Use(middlewares.ValidateSession)
	customer.GET("/list", handlers.GetCustomerListHandler)
	customer.GET("/id/:id", handlers.GetCustomerHandler)
	customer.POST("/create", handlers.CreateCustomerHandler)
	customer.PUT("/update/:id", handlers.UpdateCustomerHandler)
	customer.DELETE("/delete/:id", handlers.DeleteCustomerHandler)

	employee := rg.Group("/employee")
	// employees.Use(middlewares.ValidateSession)
	employee.GET("/list", handlers.GetEmployeeListHandler)
	employee.GET("/id/:id", handlers.GetEmployeeHandler)
	employee.POST("/create", handlers.CreateEmployeeHandler)
	employee.PUT("/update/:id", handlers.UpdateEmployeeHandler)
	employee.DELETE("/delete/:id", handlers.DeleteEmployeeHandler)

	employeeType := rg.Group("/employee-type")
	// employees.Use(middlewares.ValidateSession)
	employeeType.GET("/list", handlers.GetEmployeeTypeListHandler)
	employeeType.GET("/id/:id", handlers.GetEmployeeTypeHandler)
	employeeType.POST("/create", handlers.CreateEmployeeTypeHandler)
	employeeType.PUT("/update/:id", handlers.UpdateEmployeeTypeHandler)
	employeeType.DELETE("/delete/:id", handlers.DeleteEmployeeTypeHandler)

	deliveryLocation := rg.Group("/delivery-location")
	// employees.Use(middlewares.ValidateSession)
	deliveryLocation.GET("/list", handlers.GetDeliveryLocationListHandler)
	deliveryLocation.GET("/id/:id", handlers.GetDeliveryLocationHandler)
	deliveryLocation.POST("/create", handlers.CreateDeliveryLocationHandler)
	deliveryLocation.PUT("/update/:id", handlers.UpdateDeliveryLocationHandler)
	deliveryLocation.DELETE("/delete/:id", handlers.DeleteDeliveryLocationHandler)

}
