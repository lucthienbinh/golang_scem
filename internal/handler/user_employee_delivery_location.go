package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	"gopkg.in/validator.v2"
)

// -------------------- DELIVERY LOCATION HANDLER FUNTION --------------------

// GetDeliveryLocationListHandler in database
func GetDeliveryLocationListHandler(c *gin.Context) {
	deliveryLocations := []model.DeliveryLocation{}
	db.Order("id asc").Find(&deliveryLocations)
	c.JSON(http.StatusOK, gin.H{"delivery_location_list": &deliveryLocations})
	return
}

func getDeliveryLocationOrNotFound(c *gin.Context) (*model.DeliveryLocation, error) {
	deliveryLocation := &model.DeliveryLocation{}
	if err := db.First(deliveryLocation, c.Param("id")).Error; err != nil {
		return deliveryLocation, err
	}
	return deliveryLocation, nil
}

// GetDeliveryLocationHandler in database
func GetDeliveryLocationHandler(c *gin.Context) {
	deliveryLocation, err := getDeliveryLocationOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"delivery_location_info": &deliveryLocation})
	return
}

// CreateDeliveryLocationHandler in database
func CreateDeliveryLocationHandler(c *gin.Context) {
	deliveryLocation := &model.DeliveryLocation{}
	if err := c.ShouldBindJSON(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.Create(&deliveryLocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "A delivery location has been created!"})
	return
}

// UpdateDeliveryLocationHandler in database
func UpdateDeliveryLocationHandler(c *gin.Context) {
	deliveryLocation, err := getDeliveryLocationOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&deliveryLocation); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	deliveryLocation.ID = getIDFromParam(c)
	if err = db.Model(&deliveryLocation).Updates(&deliveryLocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteDeliveryLocationHandler in database
func DeleteDeliveryLocationHandler(c *gin.Context) {
	if _, err := getDeliveryLocationOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.DeliveryLocation{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
