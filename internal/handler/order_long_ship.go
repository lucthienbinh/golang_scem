package handler

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucthienbinh/golang_scem/internal/model"
	CommonService "github.com/lucthienbinh/golang_scem/internal/service/common"
	CommonMessage "github.com/lucthienbinh/golang_scem/internal/service/common_message"
	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/sync/errgroup"
	"gopkg.in/validator.v2"
)

// -------------------- LONG SHIP HANDLER FUNTION --------------------

// GetLongShipListHandler in database
func GetLongShipListHandler(c *gin.Context) {
	queryStrings := c.Request.URL.Query()
	sortByCondition := queryStrings.Get("sortByCondition")

	type APILongShipList struct {
		ID                       uint   `json:"id"`
		TransportTypeID          uint   `json:"transport_type_id"`
		LicensePlate             string `json:"license_plate"`
		EstimatedTimeOfDeparture int64  `json:"estimated_time_of_departure"`
		EstimatedTimeOfArrival   int64  `json:"estimated_time_of_arrival"`
		PackageLoaded            bool   `json:"package_loaded"`
		Finished                 bool   `json:"finished"`
	}
	longShips := []APILongShipList{}
	longShips2 := []APILongShipList{}
	db.Model(&model.LongShip{}).Order("id asc").Find(&longShips)
	if sortByCondition == "available" {
		db.Model(&model.LongShip{}).Order("id asc").Find(&longShips2, "estimated_time_of_departure > ?", time.Now().Unix())
	} else if sortByCondition == "ready" {
		db.Model(&model.LongShip{}).Order("id asc").Find(&longShips2, "estimated_time_of_departure < ? AND finished is ? AND package_loaded is ?", time.Now().Unix(), false, false)
	} else if sortByCondition == "running" {
		db.Model(&model.LongShip{}).Order("id asc").Find(&longShips2, "estimated_time_of_departure < ? AND finished is ? AND package_loaded is ?", time.Now().Unix(), false, true)
	} else if sortByCondition == "finished" {
		db.Model(&model.LongShip{}).Order("id asc").Find(&longShips2, "finished is ?", true)
	} else if sortByCondition == "" {
		longShips2 = longShips
	}

	transportTypes := []model.TransportType{}
	db.Where("same_city is ?", false).Order("id asc").Find(&transportTypes)

	longShipTotal := 0
	readyTotal := 0
	availableTotal := 0
	runningTotal := 0
	finishedTotal := 0
	for i := 0; i < len(longShips); i++ {
		longShipTotal++
		timeNow := time.Now()
		beforeDeptTime := timeNow.Before(time.Unix(longShips[i].EstimatedTimeOfDeparture, 0))
		afterDeptTime := timeNow.After(time.Unix(longShips[i].EstimatedTimeOfDeparture, 0)) && (longShips[i].Finished == false) && (longShips[i].PackageLoaded == false)
		afterDeptRunTime := timeNow.After(time.Unix(longShips[i].EstimatedTimeOfDeparture, 0)) && (longShips[i].Finished == false) && (longShips[i].PackageLoaded == true)

		if longShips[i].Finished {
			finishedTotal++
		}
		if beforeDeptTime {
			availableTotal++
		}
		if afterDeptTime {
			readyTotal++
		}
		if afterDeptRunTime {
			runningTotal++
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"long_ship_list":      &longShips2,
		"transport_type_list": &transportTypes,
		"long_ship_total":     longShipTotal,
		"ready_total":         readyTotal,
		"available_total":     availableTotal,
		"running_total":       runningTotal,
		"finished_total":      finishedTotal,
	})
	return
}

func getLongShipOrNotFound(c *gin.Context) (*model.LongShip, error) {
	longShip := &model.LongShip{}
	if err := db.First(longShip, c.Param("id")).Error; err != nil {
		return longShip, err
	}
	return longShip, nil
}

func getLongShipOrNotFoundByQrcode(c *gin.Context) (*model.LongShip, error) {
	longShip := &model.LongShip{}
	if err := db.First(longShip, "ls_qr_code LIKE ?", c.Param("qrcode")).Error; err != nil {
		return longShip, err
	}
	return longShip, nil
}

// GetLongShipHandler in database
func GetLongShipHandler(c *gin.Context) {
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"long_ship_info": &longShip})
	return
}

// CreateLongShipFormData in frontend
func CreateLongShipFormData(c *gin.Context) {
	transportTypes := []model.TransportType{}
	db.Where("same_city is ?", false).Order("id asc").Find(&transportTypes)
	c.JSON(http.StatusOK, gin.H{"transport_type_list": &transportTypes})
	return
}

// CreateLongShipHandler in database
func CreateLongShipHandler(c *gin.Context) {
	longShip := &model.LongShip{}
	if err := c.ShouldBindJSON(&longShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&longShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Create long ship id base on Time
	current := uuid.New().Time()
	currentString := fmt.Sprintf("%d", current)
	rawUint, _ := strconv.ParseUint(currentString, 10, 64)
	longShip.ID = uint(rawUint / 100000000000)

	// Run concurrency
	var g errgroup.Group

	// Create QR code
	g.Go(func() error {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			return err
		}
		newQrCode := fmt.Sprintf("%x", b)
		createTime := fmt.Sprintf("%d", time.Now().Unix())
		newQrCode = createTime + "_" + newQrCode + ".jpg"
		filepath := os.Getenv("QR_CODE_FILE_PATH") + newQrCode
		if err := qrcode.WriteFile(newQrCode, qrcode.Medium, 256, filepath); err != nil {
			return err
		}
		longShip.LSQrCode = newQrCode
		if err := db.Create(longShip).Error; err != nil {
			return err
		}
		return nil
	})

	// Create workflow instance in zeebe in state machine and save to long ship workflow data
	g.Go(func() error {
		WorkflowKey, WorkflowInstanceKey, err := CommonService.CreateWorkflowLongShipInstanceHandler(longShip.ID)
		if err != nil {
			return err
		}

		longShipWorkflowData := &model.LongShipWorkflowData{}
		longShipWorkflowData.LongShipID = longShip.ID
		longShipWorkflowData.WorkflowKey = WorkflowKey
		longShipWorkflowData.WorkflowInstanceKey = WorkflowInstanceKey
		if err := db.Create(longShipWorkflowData).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"server_response": "A long ship has been created!"})
	return
}

// UpdateLongShipFormData in frontend
func UpdateLongShipFormData(c *gin.Context) {
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	transportTypes := []model.TransportType{}
	db.Where("same_city == ?", false).Order("id asc").Find(&transportTypes)
	c.JSON(http.StatusOK, gin.H{
		"long_ship_info":      &longShip,
		"transport_type_list": &transportTypes,
	})
	return
}

// UpdateLongShipHandler in database
func UpdateLongShipHandler(c *gin.Context) {
	longShip, err := getLongShipOrNotFound(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&longShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := validator.Validate(&longShip); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	longShip.ID = getIDFromParam(c)
	selectField := []string{"transport_type_id", "license_plate", "estimated_time_of_departure", "estimated_time_of_arrival"}
	if err = db.Model(&longShip).Select(selectField).Updates(longShip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateLSLoadPackageHandler in database
func UpdateLSLoadPackageHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// QRcode will replace this line
	longShip, err := getLongShipOrNotFoundByQrcode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// if longShip.PackageLoaded == true {
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }
	// Run concurrency
	var g errgroup.Group

	// Send message to sate machine (Zeebe or State Scem)
	g.Go(func() error {
		if err := CommonMessage.PublishPackageLoadedMessage(longShip.ID); err != nil {
			return err
		}
		return nil
	})

	// Send notification to FCM cloud and store in database
	g.Go(func() error {
		if err := createCustomerNotificationLongShipHandler(longShip.ID, 1); err != nil {
			return err
		}
		return nil
	})

	// Update long ship in database
	g.Go(func() error {
		longShipUpdateInfo := model.LongShip{
			CurrentLocation: "Location1",
			PackageLoaded:   true,
			EmplLoadID:      employeeID,
			LoadedTime:      time.Now().Unix(),
		}
		if err := db.Model(&longShip).Updates(longShipUpdateInfo).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateLSStartVehicleHandler in database
func UpdateLSStartVehicleHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// QRcode will replace this line
	longShip, err := getLongShipOrNotFoundByQrcode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if longShip.PackageLoaded == false || longShip.VehicleStarted == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Run concurrency
	var g errgroup.Group

	// Send message to sate machine (Zeebe or State Scem)
	g.Go(func() error {
		if err := CommonMessage.PublishVehicleStartedMessage(longShip.ID); err != nil {
			return err
		}
		return nil
	})

	// Update long ship in database
	g.Go(func() error {
		longShipUpdateInfo := model.LongShip{
			CurrentLocation: "Location2",
			VehicleStarted:  true,
			EmplDriver1ID:   employeeID,
			StartedTime:     time.Now().Unix(),
		}
		if err := db.Model(&longShip).Updates(longShipUpdateInfo).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateLSVehicleArrivedHandler in database
func UpdateLSVehicleArrivedHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// QRcode will replace this line
	longShip, err := getLongShipOrNotFoundByQrcode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if longShip.VehicleStarted == false || longShip.VehicleArrived == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Run concurrency
	var g errgroup.Group

	// Send message to sate machine (Zeebe or State Scem)
	g.Go(func() error {
		if err := CommonMessage.PublishVehicleArrivedMessage(longShip.ID); err != nil {
			return err
		}
		return nil
	})

	// Update long ship in database
	g.Go(func() error {
		longShipUpdateInfo := model.LongShip{
			CurrentLocation: "Location3",
			VehicleArrived:  true,
			EmplDriver2ID:   employeeID,
			ArrivedTime:     time.Now().Unix(),
		}
		if err := db.Model(&longShip).Updates(longShipUpdateInfo).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// UpdateLSUnloadPackageHandler in database
func UpdateLSUnloadPackageHandler(c *gin.Context, userAuthID uint) {
	employeeID, err := getEmployeeIDByUserAuthID(userAuthID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// QRcode will replace this line
	longShip, err := getLongShipOrNotFoundByQrcode(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if longShip.VehicleArrived == false || longShip.PackageUnloaded == true {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Run concurrency
	var g errgroup.Group

	// Send message to state machine (Zeebe or State Scem)
	g.Go(func() error {
		if err := CommonMessage.PublishPackageUnloadedMessage(longShip.ID); err != nil {
			return err
		}
		return nil
	})

	// Send notification to FCM cloud and store in database
	g.Go(func() error {
		if err := createCustomerNotificationLongShipHandler(longShip.ID, 2); err != nil {
			return err
		}
		return nil
	})

	// Update long ship in database
	g.Go(func() error {
		longShipUpdateInfo := model.LongShip{
			CurrentLocation: "Location4",
			PackageUnloaded: true,
			EmplUnloadID:    employeeID,
			UnloadedTime:    time.Now().Unix(),
			Finished:        true,
		}
		if err := db.Model(&longShip).Updates(longShipUpdateInfo).Error; err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been updated!"})
	return
}

// DeleteLongShipHandler in database
func DeleteLongShipHandler(c *gin.Context) {
	if _, err := getLongShipOrNotFound(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&model.LongShip{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"server_response": "Your information has been deleted!"})
	return
}
