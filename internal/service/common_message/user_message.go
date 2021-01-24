package message

import (
	"log"
	"os"

	ZBMessage "github.com/lucthienbinh/golang_scem/internal/service/zeebe/message"
)

///////////////////////////////////////// PUBLISH MESSAGE /////////////////////////////////////////

////////+++++++++++ MESSAGE SELECTOR +++++++++++////////

// PublishPaymentConfirmedMessage will select private function
func PublishPaymentConfirmedMessage(orderID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.PaymentConfirmed(orderID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	// Run without any state machine
	orderWorkflowData, err := getOrderWorkflowDataByOrderID(orderID)
	if err != nil {
		return err
	}
	if orderWorkflowData.UseLongShip == true {
		orderLongShipID, err := CreateOrderLongShip(orderID)
		if err != nil {
			return err
		}
		log.Println("Created order long ship id:", orderLongShipID)
	} else {
		orderShortShipID, err := CreateOrderShortShip(orderID)
		if err != nil {
			return err
		}
		log.Println("Created order short ship id:", orderShortShipID)
	}
	return nil
}

// PublishPackageLoadedMessage will select private function
func PublishPackageLoadedMessage(longShipID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.PackageLoaded(longShipID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}

// PublishVehicleStartedMessage will select private function
func PublishVehicleStartedMessage(longShipID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.VehicleStarted(longShipID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}

// PublishVehicleArrivedMessage will select private function
func PublishVehicleArrivedMessage(longShipID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.VehicleArrived(longShipID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}

// PublishPackageUnloadedMessage will select private function
func PublishPackageUnloadedMessage(longShipID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.PackageUnloaded(longShipID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	// Run without any state machine
	// Unloaded Pakage message -> Long Ship Finished -> Order Short Ship Service
	orderLongShips, err := GetOrderLongShipList(longShipID)
	if err != nil {
		return err
	}
	for i := 0; i < len(orderLongShips); i++ {
		orderID := orderLongShips[i].OrderID
		orderShortShipID, err := CreateOrderShortShip(orderID)
		if err != nil {
			return err
		}
		log.Println("Created order short ship id:", orderShortShipID)
	}
	return nil
}

// PublishShipperCalledMessage will select private function
func PublishShipperCalledMessage(orderID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.ShipperCalled(orderID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}

// PublishShipperReceivedMoneyMessage will select private function
func PublishShipperReceivedMoneyMessage(orderID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.ShipperReceivedMoney(orderID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}

// PublishShipperShippedMessage will select private function
func PublishShipperShippedMessage(orderID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.ShipperShipped(orderID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}

// PublishShipperConfirmedMessage will select private function
func PublishShipperConfirmedMessage(orderID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.ShipperConfirmed(orderID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}
