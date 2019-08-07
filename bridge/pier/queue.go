package pier

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/streadway/amqp"
)

const (
	connector = "queue-connector"
)

type QueueConnector struct {
	AmqpDailer      string
	HeimdallQueue   string
	BorQueue        string
	CheckpointQueue string
	cliContext      cliContext.CLIContext
	Logger          log.Logger
}

// NewQueueConnector creates a connector object which can be used to connect/send/consume bytes from queue
func NewQueueConnector(dialer string, heimdallQ string, borQ string, checkpointq string) QueueConnector {
	logger := Logger.With("module", connector)

	cliCtx := cliContext.NewCLIContext()
	cliCtx.BroadcastMode = client.BroadcastAsync
	return QueueConnector{
		AmqpDailer:      dialer,
		HeimdallQueue:   heimdallQ,
		BorQueue:        borQ,
		CheckpointQueue: checkpointq,
		cliContext:      cliCtx,
		Logger:          logger,
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
	qc.Logger.Info("Dispatched message to heimdall", "MsgType", msg.Type())
	if err != nil {
		return err
	}
	return nil
}

// ConsumeHeimdallQ consumes messages from heimdall queue
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
			qc.Logger.Debug("Sending transaction to heimdall", "TxBytes", d.Body)
			resp, err := helper.SendTendermintRequest(qc.cliContext, d.Body, helper.BroadcastAsync)
			if err != nil {
				qc.Logger.Error("Unable to send transaction to heimdall", "error", err)
			} else {
				qc.Logger.Info("Sent to heimdall", "Response", resp.String())
			}
		}
	}()

	qc.Logger.Info("Starting queue consumer")
	<-forever
	return nil
}
