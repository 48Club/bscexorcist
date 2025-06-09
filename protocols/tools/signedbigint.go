// Package tools provides utility functions for protocol event parsing.
package tools

import "math/big"

// DecodeSignedInt256 converts a 32-byte slice to a signed big.Int (two's complement).
func DecodeSignedInt256(data []byte) *big.Int {
	value := new(big.Int).SetBytes(data)

	if data[0]&0x80 != 0 {
		two256 := new(big.Int).Lsh(big.NewInt(1), 256)
		value.Sub(value, two256)
	}
	return value
}
