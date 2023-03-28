package types

// SideTxMsg tx message
type SideTxMsg interface {
	GetSideSignBytes() []byte
}
