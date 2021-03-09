package service

import (
	"github.com/RichardKnop/machinery/v1/brokers/amqp"
	"github.com/RichardKnop/machinery/v1/config"
)

type Service struct {
	broker *amqp.Broker
	config *config.Config
}

func New(config *config.Config) *Service {
	return &Service{
		broker: amqp.New(config).(*amqp.Broker),
		config: config,
	}
}
