package utils

import "math/big"

// BigIntFromBytes converts a byte slice to a big.Int.
func BigIntFromBytes(data []byte) *big.Int {
	if len(data) == 0 {
		return nil
	}
	return new(big.Int).SetBytes(data)
}

// BigIntFromPointer converts a pointer to a big.Int.
func BigIntFromPointer(data *big.Int) *big.Int {
	if data == nil {
		return nil
	}
	return new(big.Int).Set(data)
}

// DecodeSignedInt256 converts a 32-byte slice to a signed big.Int (two's complement).
func DecodeSignedInt256(data []byte) *big.Int {
	value := new(big.Int).SetBytes(data)

	if data[0]&0x80 != 0 {
		two256 := new(big.Int).Lsh(big.NewInt(1), 256)
		value.Sub(value, two256)
	}
	return value
}
