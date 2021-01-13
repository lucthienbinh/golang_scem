package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucthienbinh/golang_scem/internal/model"
	"gopkg.in/validator.v2"
)

// -------------------- TRANSPORT TYPE HANDLER FUNTION --------------------

// GetTransportTypeListHandler in database
func GetTransportTypeListHandler(c *gin.Context) {
	transportTypes := []model.TransportType{}
	db.Order("id asc").Find(&transportTypes)
	c.JSON(http.StatusOK, gin.H{"transport_type_list": &transportTypes})
	return
}

func getTransportTypeOrNotFound(c *gin.Context) (*model.TransportType, error) {
	transportType := &model.TransportType{}
	if err := db.First(transportType, c.Param("id")).Error; err != nil {
		return transportType, err
	}
	return transportType, nil
}

// GetTransportTypeHandler in database
func GetTransportTypeHandler(c *gin.Context) {
	transportType, err := getTransportTypeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transport_type_info": &transportType})
	return
}

// CreateTransportTypeHandler in database
func CreateTransportTypeHandler(c *gin.Context) {
	transportType := &model.TransportType{}
	if err := c.ShouldBindJSON(&transportType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&transportType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.Create(transportType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"server_response": "An transport type has been created!"})
	return
}

// UpdateTransportTypeHandler in database
func UpdateTransportTypeHandler(c *gin.Context) {
	transportType, err := getTransportTypeOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&transportType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&transportType); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	transportType.ID = getIDFromParam(c)
	updateValue := map[string]interface{}{
		"same_city":               transportType.SameCity,
		"location_one":            transportType.LocationOne,
		"location_two":            transportType.LocationTwo,
		"bus_station_from":        transportType.BusStationFrom,
		"bus_station_to":          transportType.BusStationTo,
		"long_ship_duration":      transportType.LongShipDuration,
		"long_ship_price":         transportType.LongShipPrice,
		"short_ship_price_per_km": transportType.ShortShipPricePerKm,
	}
	if err = db.Model(&transportType).Updates(updateValue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteTransportTypeHandler in database
func DeleteTransportTypeHandler(c *gin.Context) {
	if _, err := getTransportTypeOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.TransportType{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
