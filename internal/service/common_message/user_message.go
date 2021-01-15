package message

import (
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
	return nil
}

// PublishPackageLoadedMessage will select private function
func PublishPackageLoadedMessage(orderID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.PackageLoaded(orderID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}

// PublishVehicleStartedMessage will select private function
func PublishVehicleStartedMessage(orderID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.VehicleStarted(orderID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}

// PublishVehicleArrivedMessage will select private function
func PublishVehicleArrivedMessage(orderID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.VehicleArrived(orderID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
	}
	return nil
}

// PublishPackageUnloadedMessage will select private function
func PublishPackageUnloadedMessage(orderID uint) error {
	if os.Getenv("STATE_SERVICE") == "1" {
		return ZBMessage.PackageUnloaded(orderID)
	}
	if os.Getenv("STATE_SERVICE") == "2" {
		return nil
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
