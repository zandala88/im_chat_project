package mq

import (
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"im/model"
	"im/public/mq"
)

const (
	MessageQueue        = "message.queue"
	MessageRoutingKey   = "message.routing.key"
	MessageExchangeName = "message.exchange.name"
)

var (
	MessageMQ *mq.Conn
)

func InitMessageMQ(url string) {
	MessageMQ = mq.InitRabbitMQ(url, MessageCreateHandler, MessageQueue, MessageRoutingKey, MessageExchangeName)
}

func MessageCreateHandler(d rabbitmq.Delivery) rabbitmq.Action {
	messageModels := model.ProtoMarshalToMessage(d.Body)
	if messageModels == nil {
		zap.S().Debug("[MessageCreateHandler] messageModels 为空")
		return rabbitmq.NackDiscard
	}
	err := model.CreateMessage(messageModels...)
	if err != nil {
		zap.S().Error("[MessageCreateHandler] [CreateMessage] [err] = ", err)
		return rabbitmq.NackDiscard
	}

	zap.S().Debug("[MessageCreateHandler] [CreateMessage] [messageModels] = ", messageModels)
	return rabbitmq.Ack
}
