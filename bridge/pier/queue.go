package pier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/core/types"
	"github.com/streadway/amqp"
	"github.com/tendermint/tendermint/libs/log"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	connector = "queue-connector"

	// exchanges
	broadcastExchange = "bridge.exchange.broadcast"
	// heimdall queue
	heimdallBroadcastQueue = "bridge.queue.heimdall"
	// bor queue
	borBroadcastQueue = "bridge.queue.bor"

	// heimdall routing key
	heimdallBroadcastRoute = "bridge.route.heimdall"
	// bor routing key
	borBroadcastRoute = "bridge.route.bor"
)

// QueueConnector queue connector
type QueueConnector struct {
	// URL for connecting to AMQP
	connection *amqp.Connection
	// create a channel
	channel *amqp.Channel
	// tx encoder
	cliCtx cliContext.CLIContext
	// logger
	logger log.Logger
}

// NewQueueConnector creates a connector object which can be used to connect/send/consume bytes from queue
func NewQueueConnector(cdc *codec.Codec, dialer string) *QueueConnector {
	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	// amqp dialer
	conn, err := amqp.Dial(dialer)
	if err != nil {
		panic(err)
	}

	// initialize exchange
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	// queue connector
	connector := QueueConnector{
		connection: conn,
		channel:    channel,
		cliCtx:     cliCtx,
		// create logger
		logger: Logger.With("module", "queue-connector"),
	}

	// connector
	return &connector
}

// Start connector
func (qc *QueueConnector) Start() error {
	// exchange declare
	if err := qc.channel.ExchangeDeclare(
		broadcastExchange, // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	); err != nil {
		return err
	}

	// Heimdall

	// queue declare
	if _, err := qc.channel.QueueDeclare(
		heimdallBroadcastQueue, // name
		true,                   // durable
		false,                  // delete when usused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	); err != nil {
		return err
	}

	// bind queue
	if err := qc.channel.QueueBind(
		heimdallBroadcastQueue, // queue name
		heimdallBroadcastRoute, // routing key
		broadcastExchange,      // exchange
		false,
		nil,
	); err != nil {
		return err
	}

	// start consuming
	msgs, err := qc.channel.Consume(
		heimdallBroadcastQueue, // queue
		heimdallBroadcastQueue, // consumer  -- consumer identifier
		false,                  // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	if err != nil {
		return err
	}
	// process heimdall broadcast messages
	go qc.handleHeimdallBroadcastMsgs(msgs)

	// Bor

	// queue declare
	if _, err := qc.channel.QueueDeclare(
		borBroadcastQueue, // name
		true,              // durable
		false,             // delete when usused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	); err != nil {
		return err
	}

	// bind queue
	if err := qc.channel.QueueBind(
		borBroadcastQueue, // queue name
		borBroadcastRoute, // routing key
		broadcastExchange, // exchange
		false,
		nil,
	); err != nil {
		return err
	}

	// start consuming
	msgs, err = qc.channel.Consume(
		borBroadcastQueue, // queue
		borBroadcastQueue, // consumer  -- consumer identifier
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return err
	}

	// process bor broadcast messages
	go qc.handleBorBroadcastMsgs(msgs)

	return nil
}

// Stop connector
func (qc *QueueConnector) Stop() {
	// close channel & connection
	qc.channel.Close()
	qc.connection.Close()
}

//
// Publish
//

// BroadcastToHeimdall broadcasts to heimdall
func (qc *QueueConnector) BroadcastToHeimdall(msg sdk.Msg) error {
	data, err := qc.cliCtx.Codec.MarshalJSON(msg)
	if err != nil {
		return err
	}

	return qc.BroadcastBytesToHeimdall(data)
}

// BroadcastBytesToHeimdall broadcasts bytes to heimdall
func (qc *QueueConnector) BroadcastBytesToHeimdall(data []byte) error {
	if err := qc.channel.Publish(
		broadcastExchange,      // exchange
		heimdallBroadcastRoute, // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		}); err != nil {
		return err
	}

	return nil
}

// BroadcastToBor broadcasts to bor
func (qc *QueueConnector) BroadcastToBor(data []byte) error {
	if err := qc.channel.Publish(
		broadcastExchange, // exchange
		borBroadcastRoute, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		}); err != nil {
		return err
	}

	return nil
}

//
// Consume
//

func (qc *QueueConnector) handleHeimdallBroadcastMsgs(amqpMsgs <-chan amqp.Delivery) {
	// tx encoder
	txEncoder := helper.GetTxEncoder(qc.cliCtx.Codec)
	// chain id
	chainID := helper.GetGenesisDoc().ChainID
	// current address
	address := hmTypes.BytesToHeimdallAddress(helper.GetAddress())
	// fetch from APIs
	var account authTypes.Account
	response, err := FetchFromAPI(qc.cliCtx, GetHeimdallServerEndpoint(fmt.Sprintf(AccountDetailsURL, address)))
	if err != nil {
		qc.logger.Error("Error fetching account from rest-api", "url", GetHeimdallServerEndpoint(fmt.Sprintf(AccountDetailsURL, address)))
		panic("Error connecting to rest-server, please start server before bridge")
	}

	// get proposer from response
	if err := qc.cliCtx.Codec.UnmarshalJSON(response.Result, &account); err != nil && len(response.Result) != 0 {
		panic(err)
	}

	// get account number and sequence
	accNum := account.GetAccountNumber()
	accSeq := account.GetSequence()

	// handler
	handler := func(amqpMsg amqp.Delivery) bool {
		var msg sdk.Msg
		if err := qc.cliCtx.Codec.UnmarshalJSON(amqpMsg.Body, &msg); err != nil {
			amqpMsg.Reject(false)
			qc.logger.Error("Error while broadcasting the heimdall transaction", "error", err)
			return false
		}

		txBldr := authTypes.NewTxBuilderFromCLI().
			WithTxEncoder(txEncoder).
			WithAccountNumber(accNum).
			WithSequence(accSeq).
			WithChainID(chainID)
		if _, err := helper.BuildAndBroadcastMsgs(qc.cliCtx, txBldr, []sdk.Msg{msg}); err != nil {
			amqpMsg.Reject(false)
			qc.logger.Error("Error while broadcasting the heimdall transaction", "error", err)
			return false
		}

		// send ack
		amqpMsg.Ack(false)

		// increment account sequence
		accSeq = accSeq + 1

		return true
	}

	// handle all amqp messages
	for amqpMsg := range amqpMsgs {
		handler(amqpMsg)
	}
}

func (qc *QueueConnector) handleBorBroadcastMsgs(amqpMsgs <-chan amqp.Delivery) {
	maticClient := helper.GetMaticClient()

	// handler
	handler := func(amqpMsg amqp.Delivery) bool {
		var msg ethereum.CallMsg
		if err := json.Unmarshal(amqpMsg.Body, &msg); err != nil {
			amqpMsg.Reject(false)
			qc.logger.Error("Error while parsing the transaction from queue", "error", err)
			return false
		}

		// get auth
		auth, err := helper.GenerateAuthObj(maticClient, *msg.To, msg.Data)
		if err != nil {
			amqpMsg.Reject(false)
			qc.logger.Error("Error while fetching the transaction param details", "error", err)
			return false
		}

		// Create the transaction, sign it and schedule it for execution
		rawTx := types.NewTransaction(auth.Nonce.Uint64(), *msg.To, msg.Value, auth.GasLimit, auth.GasPrice, msg.Data)
		// signer
		signedTx, err := auth.Signer(types.HomesteadSigner{}, auth.From, rawTx)
		if err != nil {
			amqpMsg.Reject(false)
			qc.logger.Error("Error while signing the transaction", "error", err)
			return false
		}

		qc.logger.Debug("Sending transaction to bor", "TxHash", signedTx.Hash())

		// broadcast transaction
		if err := maticClient.SendTransaction(context.Background(), signedTx); err != nil {
			amqpMsg.Reject(false)
			qc.logger.Error("Error while broadcasting the transaction", "error", err)
			return false
		}

		// send ack
		amqpMsg.Ack(false)

		// amqp msg
		return true
	}

	// handle all amqp messages
	for amqpMsg := range amqpMsgs {
		handler(amqpMsg)
	}
}
