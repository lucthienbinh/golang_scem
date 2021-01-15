package message

import (
	"context"
	"fmt"
	"time"

	CommonService "github.com/lucthienbinh/golang_scem/internal/service/common"
)

// LongShipFinished send mesage to zeebe engine
func LongShipFinished(longShipID uint) error {
	ctx := context.Background()
	orderLongShips, err := CommonService.GetOrderLongShipList(longShipID)
	if err != nil {
		return err
	}
	for i := 0; i < len(orderLongShips); i++ {
		orderID := orderLongShips[i].OrderID
		_, err := zbClient.NewPublishMessageCommand().MessageName("LongShipFinished").CorrelationKey(fmt.Sprint(orderID)).TimeToLive(1 * time.Minute).Send(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
