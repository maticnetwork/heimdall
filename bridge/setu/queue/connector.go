package queue

import (
	"github.com/streadway/amqp"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
)

type QueueConnector struct {
	connection        *amqp.Connection
	broadcastExchange string
	logger            log.Logger
	Server            *machinery.Server
}

func NewQueueConnector(dialer string) *QueueConnector {
	// amqp dialer
	conn, err := amqp.Dial(dialer)
	if err != nil {
		panic(err)
	}

	var cnf = &config.Config{
		Broker:        dialer,
		DefaultQueue:  "machinery_tasks",
		ResultBackend: dialer,
		AMQP: &config.AMQPConfig{
			Exchange:     "machinery_exchange",
			ExchangeType: "direct",
			BindingKey:   "machinery_task",
		},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		// do something with the error
	}

	// queue connector
	connector := QueueConnector{
		connection:        conn,
		broadcastExchange: BroadcastExchange,
		logger:            util.Logger().With("module", Connector),
		Server:            server,
	}

	// connector
	return &connector
}

// StartWorker - starts worker to process registered tasks
func (qc *QueueConnector) StartWorker() {
	worker := qc.Server.NewWorker("invoke-processor", 10)
	qc.logger.Info("Starting worker")
	errors := make(chan error)
	worker.LaunchAsync(errors)
}

// // InitializeQueues initiates multiple queues and exchange
// func (qc *QueueConnector) InitializeQueues() error {
// 	// initialize exchange
// 	channel, err := qc.connection.Channel()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// exchange declare
// 	if err := channel.ExchangeDeclare(
// 		qc.broadcastExchange, // name
// 		"topic",              // type
// 		true,                 // durable
// 		false,                // auto-deleted
// 		false,                // internal
// 		false,                // no-wait
// 		nil,                  // arguments
// 	); err != nil {
// 		return err
// 	}

// 	qc.logger.Debug("AMQP exchange declared", "exchange", qc.broadcastExchange)

// 	qc.InitializeQueue(channel, CheckpointQueueName, CheckpointQueueRoute)
// 	qc.InitializeQueue(channel, StakingQueueName, StakingQueueRoute)
// 	qc.InitializeQueue(channel, FeeQueueName, FeeQueueRoute)
// 	qc.InitializeQueue(channel, SpanQueueName, SpanQueueRoute)
// 	qc.InitializeQueue(channel, ClerkQueueName, ClerkQueueRoute)
// 	return nil
// }

// // InitializeQueue initialize individual queue
// func (qc *QueueConnector) InitializeQueue(channel *amqp.Channel, queueName string, queueRoute string) error {
// 	// queue declare
// 	if _, err := channel.QueueDeclare(
// 		queueName, // name
// 		true,      // durable
// 		false,     // delete when usused
// 		false,     // exclusive
// 		false,     // no-wait
// 		nil,       // arguments
// 	); err != nil {
// 		return err
// 	}

// 	// bind queue
// 	if err := channel.QueueBind(
// 		queueName,            // queue name
// 		queueRoute,           // routing key
// 		qc.broadcastExchange, // exchange
// 		false,
// 		nil,
// 	); err != nil {
// 		return err
// 	}

// 	qc.logger.Debug("AMQP queue declared", "queue", queueName, "route", queueRoute)

// 	return nil
// }

// // PublishMsg publishes messages to queue
// func (qc *QueueConnector) PublishMsg(data []byte, route string, appID string, msgType string) error {
// 	// initialize exchange
// 	channel, err := qc.connection.Channel()
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer channel.Close()

// 	if err := channel.Publish(
// 		qc.broadcastExchange, // exchange
// 		route,                // routing key
// 		false,                // mandatory
// 		false,                // immediate
// 		amqp.Publishing{
// 			AppId:       appID,
// 			Type:        msgType,
// 			ContentType: "text/plain",
// 			Body:        data,
// 		}); err != nil {
// 		return err
// 	}

// 	qc.logger.Debug("Published message to queue", "appID", appID, "route", route, "msgType", msgType)
// 	return nil
// }

// // ConsumeMsg consume messages
// func (qc *QueueConnector) ConsumeMsg(queueName string) (<-chan amqp.Delivery, error) {
// 	// initialize exchange
// 	channel, err := qc.connection.Channel()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// start consuming
// 	msgs, err := channel.Consume(
// 		queueName, // queue
// 		queueName, // consumer  -- consumer identifier
// 		false,     // auto-ack
// 		false,     // exclusive
// 		false,     // no-local
// 		false,     // no-wait
// 		nil,       // args
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return msgs, nil
// }
