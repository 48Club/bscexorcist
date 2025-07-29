package types

import (
	"github.com/ethereum/go-ethereum/common"
)

// Addresses is a 32-byte array alias
type Addresses [32]byte

// AddressesB20 creates an Addresses from a 20-byte address (padded to 32 bytes)
func AddressesB20(addr common.Address) Addresses {
	var result [32]byte
	copy(result[12:], addr[:])
	return result
}

// AddressesB32 creates an Addresses from a 32-byte address
func AddressesB32(addr common.Hash) Addresses {
	return [32]byte(addr)
}

// IsB20 checks if the address is a 20-byte address (padded to 32 bytes)
func (a Addresses) IsB20() bool {
	for _, b := range a[:12] {
		if b != 0 {
			return false
		}
	}
	return true
}

// ToB20 converts the address to a common.Address (20 bytes)
func (a Addresses) ToB20() common.Address {
	var addr common.Address
	copy(addr[:], a[12:])
	return addr
}

// ToB32 converts the address to a common.Hash (already 32 bytes)
func (a Addresses) ToB32() common.Hash {
	if a.IsB20() {
		var result [32]byte
		copy(result[12:], a[12:])
		return common.Hash(result)
	}
	return common.Hash(a)
}

// String returns the string representation of the address
func (a Addresses) String() string {
	if a.IsB20() {
		return a.ToB20().Hex()
	}
	return a.ToB32().Hex()
}
