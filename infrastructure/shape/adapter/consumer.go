package adapter

import (
	"errors"

	"github.com/Space-Core-Dev/Mistaken.Utils/result"
	"github.com/rabbitmq/amqp091-go"
)

type ConsumerParams struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp091.Table
	QosParams
}

type QosParams struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}

func (c ConsumerParams) HasQueue() bool {
	return c.Queue != ""
}

func (q QosParams) HasQos() bool {
	return q.PrefetchCount != 0
}

func PrepareConsumer(ch *amqp091.Channel, consumerParams ConsumerParams) result.Result[<-chan amqp091.Delivery] {
	if !consumerParams.HasQueue() {
		return result.Fail[<-chan amqp091.Delivery](errors.New("queue name shouldn't be empty"))
	}

	if consumerParams.QosParams.HasQos() {
		p := consumerParams.QosParams
		err := ch.Qos(p.PrefetchCount, p.PrefetchSize, p.Global)
		if err != nil {
			return result.Fail[<-chan amqp091.Delivery](err)
		}
	}

	consumer, err := ch.Consume(consumerParams.Queue,
		consumerParams.Consumer,
		consumerParams.AutoAck,
		consumerParams.Exclusive,
		consumerParams.NoLocal,
		consumerParams.NoWait,
		consumerParams.Args)
	if err != nil {
		return result.Fail[<-chan amqp091.Delivery](err)
	}

	return result.OK(consumer)
}
