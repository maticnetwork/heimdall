package queue

import (
	"os"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/streadway/amqp"
)

type QueueConnector struct {
	connection        *amqp.Connection
	broadcastExchange string
	logger            log.Logger
}

// Global logger for bridge
var Logger log.Logger

func init() {
	Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

func NewQueueConnector(dialer string) *QueueConnector {

	// amqp dialer
	conn, err := amqp.Dial(dialer)
	if err != nil {
		panic(err)
	}

	// queue connector
	connector := QueueConnector{
		connection:        conn,
		broadcastExchange: "broadcastexchange",
		logger:            Logger.With("module", "queue-connector"),
	}

	// connector
	return &connector
}

func (qc *QueueConnector) InitializeQueue() error {
	// initialize exchange
	channel, err := qc.connection.Channel()
	if err != nil {
		panic(err)
	}

	// exchange declare
	if err := channel.ExchangeDeclare(
		qc.broadcastExchange, // name
		"topic",              // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	); err != nil {
		return err
	}

	qc.logger.Info("Exchange Declared")

	// queue declare
	if _, err := channel.QueueDeclare(
		"test-queue", // name
		true,         // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		return err
	}

	qc.logger.Info("Queue Declared")

	// bind queue
	if err := channel.QueueBind(
		"test-queue",         // queue name
		"test",               // routing key
		qc.broadcastExchange, // exchange
		false,
		nil,
	); err != nil {
		return err
	}

	qc.logger.Info("Queue Bind")

	return nil
}

// PublishBytes publishes messages to queue
func (qc *QueueConnector) PublishMsg(data []byte, route string) error {
	// initialize exchange
	channel, err := qc.connection.Channel()
	if err != nil {
		panic(err)
	}

	if err := channel.Publish(
		qc.broadcastExchange, // exchange
		route,                // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		}); err != nil {
		return err
	}

	qc.logger.Info("published message to queue")
	return nil
}

func (qc *QueueConnector) ConsumeMsg(queue string) (<-chan amqp.Delivery, error) {
	// initialize exchange
	channel, err := qc.connection.Channel()
	if err != nil {
		panic(err)
	}
	// start consuming
	msgs, err := channel.Consume(
		queue, // queue
		queue, // consumer  -- consumer identifier
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
