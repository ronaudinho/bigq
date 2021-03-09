package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ronaudinho/bigq/internal/model"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/streadway/amqp"
)

func (s *Service) RecvArgo(task *model.Task) error {
	server, err := machinery.NewServer(s.config)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(task.Payload)
	if err != nil {
		return err
	}

	args := []tasks.Arg{
		tasks.Arg{
			Name:  "payload",
			Type:  "string",
			Value: string(payload),
		},
	}

	var rrk string
	if task.RoutingKey != "" {
		rrk = task.RoutingKey + "_result"
	}
	res, err := server.SendTask(&tasks.Signature{
		Name:       task.Name,
		RoutingKey: task.RoutingKey,
		Args:       args,
		// NOTE here we send the result of the queued task to another queue
		// so that it can be picked up whenever consumers want
		// as the bulk of the work is done by machinery on initial task
		// we can simply connect to the queue here without worker
		OnSuccess: []*tasks.Signature{
			&tasks.Signature{
				Name:       task.Name + "_result",
				RoutingKey: rrk,
			},
		},
	})
	if err != nil {
		return err
	}
	// NOTE we do not want to block the HTTP call, simply get information if task is queued successfully
	state := res.GetState()
	task.ID = state.TaskUUID
	task.Status = state.State
	if state.State == tasks.StateFailure {
		task.Error = state.Error
		return err
	}
	return nil
}

func (s *Service) SendArgo(queueName string) (*tasks.Signature, error) {
	conn, channel, queue, _, _, err := s.broker.Connect(
		s.config.Broker,
		s.config.MultipleBrokerSeparator,
		s.config.TLSConfig,
		s.config.AMQP.Exchange,     // exchange name
		s.config.AMQP.ExchangeType, // exchange type
		queueName,                  // queue name
		true,                       // queue durable
		false,                      // queue delete when unused
		s.config.AMQP.BindingKey,   // queue binding key
		nil,                        // exchange declare args
		amqp.Table(s.config.AMQP.QueueDeclareArgs), // queue declare args
		amqp.Table(s.config.AMQP.QueueBindingArgs), // queue binding args
	)
	if err != nil {
		return nil, err
	}

	delivery, ok, err := channel.Get(queue.Name, false)
	// nothing to do here
	if !ok {
		log.Println("nothing queued")
		return nil, nil
	}
	for err != nil {
		if !s.broker.GetRetry() {
			return nil, fmt.Errorf("Queue consume error: %s", err)
		}
		delivery, ok, err = channel.Get(queue.Name, true)
	}
	if len(delivery.Body) == 0 {
		delivery.Nack(true, false)                          // multiple, requeue
		return nil, fmt.Errorf("Received an empty message") // RabbitMQ down?
	}

	var multiple, requeue = false, false
	signature := &tasks.Signature{}
	decoder := json.NewDecoder(bytes.NewReader(delivery.Body))
	decoder.UseNumber()
	if err := decoder.Decode(signature); err != nil {
		delivery.Nack(multiple, requeue)
		return nil, fmt.Errorf("unmarshal error: %v: %v", delivery.Body, err)
	}

	log.Printf("Received new message: %s", delivery.Body)
	err = delivery.Ack(multiple)
	if err != nil {
		return signature, err
	}
	err = s.broker.Close(channel, conn)
	return signature, err
}
