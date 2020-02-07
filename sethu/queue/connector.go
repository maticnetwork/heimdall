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
		broadcastExchange: BroadcastExchange,
		logger:            Logger.With("module", Connector),
	}

	// connector
	return &connector
}

func (qc *QueueConnector) InitializeQueues() error {
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

	qc.InitializeQueue(channel, CheckpointQueueName, CheckpointQueueRoute)
	qc.InitializeQueue(channel, StakingQueueName, StakingQueueRoute)
	qc.InitializeQueue(channel, SpanQueueName, SpanQueueRoute)
	qc.InitializeQueue(channel, ClerkQueueName, ClerkQueueRoute)
	return nil
}

func (qc *QueueConnector) InitializeQueue(channel *amqp.Channel, queueName string, queueRoute string) error {
	// queue declare
	if _, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	); err != nil {
		return err
	}

	qc.logger.Info("Queue Declared", "queuename", queueName)

	// bind queue
	if err := channel.QueueBind(
		queueName,            // queue name
		queueRoute,           // routing key
		qc.broadcastExchange, // exchange
		false,
		nil,
	); err != nil {
		return err
	}

	qc.logger.Info("Queue Bind", "queuename", queueName, "queueroute", queueRoute)
	return nil
}

// PublishBytes publishes messages to queue
func (qc *QueueConnector) PublishMsg(data []byte, route string, appId string) error {
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
			AppId:       appId,
			ContentType: "text/plain",
			Body:        data,
		}); err != nil {
		return err
	}

	qc.logger.Info("published message to queue")
	return nil
}

func (qc *QueueConnector) ConsumeMsg(queueName string) (<-chan amqp.Delivery, error) {
	// initialize exchange
	channel, err := qc.connection.Channel()
	if err != nil {
		panic(err)
	}
	// start consuming
	msgs, err := channel.Consume(
		queueName, // queue
		queueName, // consumer  -- consumer identifier
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
