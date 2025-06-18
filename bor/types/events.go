package types

// staking module event types
const (
	EventTypeProposeSpan = "propose-span"

	AttributeKeySuccess          = "success"
	AttributeKeySpanID           = "span-id"
	AttributeKeySpanStartBlock   = "start-block"
	AttributeKeySpanEndBlock     = "end-block"
	AttributesKeyLatestSpanId    = "latest-span-id"
	AttributesKeyLatestBorSpanId = "latest-bor-span-id"

	AttributeValueCategory = ModuleName
)
