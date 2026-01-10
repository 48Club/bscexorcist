// Package pancakev4 provides swap event parsing for PancakeSwap V4 pools.
package pancakev4

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// PancakeV4Swap implements SwapEvent for PancakeSwap V4-style pools.
type PancakeV4Swap struct {
	poolID  [32]byte // poolID is a 32-byte identifier for the pool
	amount0 *big.Int // int128 amount0
	amount1 *big.Int // int128 amount1
}

// PairID returns a pseudo-address derived from the first 20 bytes of poolID.
func (s *PancakeV4Swap) PairID() common.Address {
	return common.BytesToAddress(s.poolID[:20])
}

// IsToken0To1 returns true if the swap direction is token0 -> token1.
func (s *PancakeV4Swap) IsToken0To1() bool {
	// amount0 > 0 means token0 is entering the pool (sell token0, buy token1)
	return s.amount0.Sign() > 0
}

// AmountIn returns the input amount for the swap.
func (s *PancakeV4Swap) AmountIn() *big.Int {
	if s.amount0.Sign() > 0 {
		return s.amount0
	}
	return new(big.Int).Neg(s.amount1)
}

// AmountOut returns the output amount for the swap.
func (s *PancakeV4Swap) AmountOut() *big.Int {
	if s.amount0.Sign() < 0 {
		return new(big.Int).Neg(s.amount0)
	}
	return s.amount1
}

// decodeSignedInt128 converts a 32-byte slice containing a signed int128 to big.Int.
// The int128 value is stored in the lower 16 bytes (right-aligned, sign-extended).
func decodeSignedInt128(data []byte) *big.Int {
	if len(data) < 32 {
		return big.NewInt(0)
	}

	// For ABI-encoded int128, the value is sign-extended to 32 bytes
	// Check the sign bit (first byte for sign extension)
	value := new(big.Int).SetBytes(data)

	// If the first byte has the high bit set, it's negative (two's complement)
	if data[0]&0x80 != 0 {
		// It's a negative number, subtract 2^256
		two256 := new(big.Int).Lsh(big.NewInt(1), 256)
		value.Sub(value, two256)
	}
	return value
}

// ParseSwap parses a PancakeSwap V4 swap log into a PancakeV4Swap struct.
// PancakeSwap V4 Swap event:
// Swap(bytes32 indexed id, address indexed sender, int128 amount0, int128 amount1, uint160 sqrtPriceX96, uint128 liquidity, int24 tick, uint24 fee, uint16 protocolFee)
// Returns nil if the log is not a valid swap event.
func ParseSwap(log *types.Log) *PancakeV4Swap {
	// Need 3 topics: event signature, id, sender
	// Data needs at least 7 fields (amount0, amount1, sqrtPriceX96, liquidity, tick, fee, protocolFee) = 7 * 32 = 224 bytes
	if len(log.Topics) != 3 || len(log.Data) < 224 {
		return nil
	}

	var poolID [32]byte
	copy(poolID[:], log.Topics[1].Bytes())

	// amount0 and amount1 are int128 but ABI-encoded as 32 bytes each (sign-extended)
	amount0 := decodeSignedInt128(log.Data[:32])
	amount1 := decodeSignedInt128(log.Data[32:64])

	return &PancakeV4Swap{
		poolID:  poolID,
		amount0: amount0,
		amount1: amount1,
	}
}
