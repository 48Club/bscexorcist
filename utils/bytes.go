package utils

// GetBytes32 returns a 32-byte slice from the input slice.
func GetBytes32(input []byte) [32]byte {
	var out [32]byte
	copy(out[:], input)
	return out
}
