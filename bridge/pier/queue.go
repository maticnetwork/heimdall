package pier

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/streadway/amqp"
)

type QueueConnector struct {
	AmqpDailer      string
	HeimdallQueue   string
	BorQueue        string
	CheckpointQueue string
}

// NewQueueConnector creates a connector object which can be used to connect/send/consume bytes from queue
func NewQueueConnector(dialer string, heimdallQ string, borQ string, checkpointq string) QueueConnector {
	return QueueConnector{
		AmqpDailer:      dialer,
		HeimdallQueue:   heimdallQ,
		BorQueue:        borQ,
		CheckpointQueue: checkpointq,
	}
}

// DispatchToBor dispatches transactions to bor
// contains deposits, state-syncs, commit-span type transactions
func (qc *QueueConnector) DispatchToBor() {

}

// DispatchToEth dispatches transactions to Ethereum chain
// contains checkpoint, validator slashing etc type transactions
func (qc *QueueConnector) DispatchToEth() {

}

// DispatchToEth dispatches transactions to Ethereum chain
// contains validator joined, validator updated etc type transactions
func (qc *QueueConnector) DispatchToHeimdall(msg sdk.Msg) error {
	conn, err := amqp.Dial(qc.AmqpDailer)
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		qc.HeimdallQueue, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}
	txBytes, err := helper.CreateTxBytes(msg)
	if err != nil {
		return err
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        txBytes,
		})
	log.Printf(" [x] Sent %s", txBytes)
	if err != nil {
		return err
	}
	return nil
}

func (qc *QueueConnector) ConsumeHeimdallQ() error {
	conn, err := amqp.Dial(qc.AmqpDailer)
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		qc.HeimdallQueue, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
