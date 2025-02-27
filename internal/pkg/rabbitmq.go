package pkg

import (
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"payment-gateway/internal/app/domain"
)

type RabbitMQ struct {
	conn                 *amqp.Connection
	channel              *amqp.Channel
	jsonTxStatusMessages chan []byte
	soapTxStatusMessages chan []byte
}

const RabbitmqExchange = "command-exchange"
const RabbitmqSoapQueue = "soap-command-queue"
const RabbitmqSoapRoutingKey = "soap-command-routing-key"
const RabbitmqJsonQueue = "json-command-queue"
const RabbitmqJsonRoutingKey = "json-command-routing-key"

func ConnectRabbitMQ(rabbitMQURL string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = channel.ExchangeDeclare(
		RabbitmqExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("could not declare exchange")
	}

	_, err = channel.QueueDeclare(
		RabbitmqSoapQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("could not declare soap queue")
	}

	_, err = channel.QueueDeclare(
		RabbitmqJsonQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("could not declare json queue")
	}

	err = channel.QueueBind(
		RabbitmqSoapQueue,
		RabbitmqSoapRoutingKey,
		RabbitmqExchange,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("could not bind soap queue to exchange")
	}

	err = channel.QueueBind(
		RabbitmqJsonQueue,
		RabbitmqJsonRoutingKey,
		RabbitmqExchange,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("could not bind json queue to exchange")
	}

	return &RabbitMQ{
		conn:                 conn,
		channel:              channel,
		jsonTxStatusMessages: make(chan []byte),
		soapTxStatusMessages: make(chan []byte),
	}, nil
}

func (r *RabbitMQ) PublishSoapData(soapData []byte) error {
	err := r.channel.Publish(
		RabbitmqExchange,
		RabbitmqSoapRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/xml",
			Body:        soapData,
		},
	)
	if err != nil {
		return domain.ErrTransactionNotPublished
	}
	return nil
}

func (r *RabbitMQ) PublishJsonData(jsonData []byte) error {
	err := r.channel.Publish(
		RabbitmqExchange,
		RabbitmqJsonRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		return domain.ErrTransactionNotPublished
	}
	return nil
}

// GetSoapMessage jsonTxStatusMessages from the queue
func (r *RabbitMQ) GetSoapMessage() chan []byte {
	messages, err := r.channel.Consume(
		RabbitmqSoapQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to consume jsonTxStatusMessages from RabbitMQ")
	}
	go func() {
		for msg := range messages {
			r.soapTxStatusMessages <- msg.Body
		}
		close(r.soapTxStatusMessages)
	}()
	return r.soapTxStatusMessages
}

// GetJsonMessage jsonTxStatusMessages from the queue
func (r *RabbitMQ) GetJsonMessage() chan []byte {
	messages, err := r.channel.Consume(
		RabbitmqJsonQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to consume jsonTxStatusMessages from RabbitMQ")
	}

	go func() {
		for msg := range messages {
			r.jsonTxStatusMessages <- msg.Body
		}
		close(r.jsonTxStatusMessages)
	}()
	return r.jsonTxStatusMessages
}

// Close RabbitMQ connection and channel
func (r *RabbitMQ) Close() {
	err := r.channel.Close()
	if err != nil {
		return
	}
	err = r.conn.Close()
	if err != nil {
		return
	}
}
