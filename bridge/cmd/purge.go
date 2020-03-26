package cmd

import (
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

// purgeCmd represents the reset of queue
var purgeCmd = &cobra.Command{
	Use:   "purge-queue",
	Short: "Reset bridge queue tasks",
	Run: func(cmd *cobra.Command, args []string) {		
		// purge Queue
		purgeQueue()
	},
}


func purgeQueue() {	
	dialer := helper.GetConfig().AmqpURL

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

	channel.QueuePurge(queue.QueueName, false)
}



func init() {
	rootCmd.AddCommand(purgeCmd)
}
