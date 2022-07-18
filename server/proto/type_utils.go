package proto

import (
	"encoding/binary"
)

func ConvertH160toAddress(h160 *H160) [20]byte {
	var addr [20]byte

	binary.BigEndian.PutUint64(addr[0:], h160.Hi.Hi)
	binary.BigEndian.PutUint64(addr[8:], h160.Hi.Lo)
	binary.BigEndian.PutUint32(addr[16:], h160.Lo)

	return addr
}

func ConvertAddressToH160(addr [20]byte) *H160 {
	return &H160{
		Lo: binary.BigEndian.Uint32(addr[16:]),
		Hi: &H128{Lo: binary.BigEndian.Uint64(addr[8:]), Hi: binary.BigEndian.Uint64(addr[0:])},
	}
}
