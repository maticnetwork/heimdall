package queue

const (
	Connector = "queue-connector"

	BroadcastExchange = "sethu.exchange.broadcast"

	// Queue & Routes

	StakingQueueName  = "queue.name.staking"
	StakingQueueRoute = "queue.route.staking"

	SpanQueueName  = "queue.name.span"
	SpanQueueRoute = "queue.route.span"

	CheckpointQueueName  = "queue.name.checkpoint"
	CheckpointQueueRoute = "queue.route.checkpoint"

	ClerkQueueName  = "queue.name.clerk"
	ClerkQueueRoute = "queue.route.clerk"
)
