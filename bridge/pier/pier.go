package pier

const (
	chainSyncer       = "chain-syncer"
	maticCheckpointer = "matic-checkpointer"

	redisURL     = "redis-url"
	lastBlockKey = "last-block" // redis key

	defaultPollInterval      = 5 * 1000                // in milliseconds
	defaultMainPollInterval  = 5 * 1000                // in milliseconds
	defaultCheckpointLength  = 256                     // checkpoint number starts with 0, so length = defaultCheckpointLength -1
	maxCheckpointLength      = 4096                    // max blocks in one checkpoint
	defaultForcePushInterval = maxCheckpointLength * 2 // in seconds (4096 * 2 seconds)
)
