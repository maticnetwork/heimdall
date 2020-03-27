package queue

import (
	"github.com/streadway/amqp"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
)

type QueueConnector struct {
	logger log.Logger
	Server *machinery.Server
}

const (
	// machinery task queue
	QueueName = "machinery_tasks"
)

func NewQueueConnector(dialer string) *QueueConnector {
	// amqp dialer
	_, err := amqp.Dial(dialer)
	if err != nil {
		panic(err)
	}

	var cnf = &config.Config{
		Broker:        dialer,
		DefaultQueue:  QueueName,
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
		logger: util.Logger().With("module", "QueueConnector"),
		Server: server,
	}

	// connector
	return &connector
}

// StartWorker - starts worker to process registered tasks
func (qc *QueueConnector) StartWorker() {
	worker := qc.Server.NewWorker("invoke-processor", 10)
	qc.logger.Info("Starting machinery worker")
	errors := make(chan error)
	worker.LaunchAsync(errors)
}
