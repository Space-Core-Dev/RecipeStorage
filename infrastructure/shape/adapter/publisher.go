package adapter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Space-Core-Dev/Mistaken.Utils/result"
	"github.com/rabbitmq/amqp091-go"
)

type deliveryMode uint8 // Should be the Transient constant or the Persistent constant. Other values mean error

const (
	_ deliveryMode = iota
	Transient
	Persistent
)

func Publish(ctx context.Context, ch *amqp091.Channel, deliveryMode deliveryMode, exchange, key string, event any) result.Result[result.Void] {
	if deliveryMode < Transient || deliveryMode > Persistent {
		return result.SimpleFail(fmt.Errorf("a wrong Delivery Mode parameter %d", deliveryMode))
	}

	b, err := json.Marshal(event)
	if err != nil {
		return result.SimpleFail(err)
	}

	err = ch.PublishWithContext(ctx, exchange,
		key,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         b,
			DeliveryMode: uint8(deliveryMode),
		})
	if err != nil {
		return result.SimpleFail(err)
	}

	return result.SimpleOk()
}
