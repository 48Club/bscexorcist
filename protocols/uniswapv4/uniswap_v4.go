// Package uniswapv4 provides swap event parsing for Uniswap V4 and compatible protocols.
package uniswapv4

import (
	"github.com/48Club/bscexorcist/protocols/tools"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SwapEventSignature for Uniswap V4
var SwapEventSignature = common.HexToHash("0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f")

// V4Swap implements SwapEvent for Uniswap V4-style pools.
type V4Swap struct {
	poolID  [32]byte // poolID is a 32-byte identifier for the pool, used as a unique pool identifier
	amount0 *big.Int
	amount1 *big.Int
}

// PairID returns a pseudo-address derived from the first 20 bytes of poolID.
func (s *V4Swap) PairID() common.Address {
	// Use the first 20 bytes of poolID as a virtual address
	return common.BytesToAddress(s.poolID[:20])
}

// IsToken0To1 returns true if the swap direction is token0 -> token1.
func (s *V4Swap) IsToken0To1() bool {
	// amount0 > 0 means token0 is entering the pool (sell token0, buy token1)
	return s.amount0.Sign() > 0
}

// AmountIn returns the input amount for the swap.
func (s *V4Swap) AmountIn() *big.Int {
	if s.amount0.Sign() > 0 {
		return s.amount0
	}
	return new(big.Int).Neg(s.amount1)
}

// AmountOut returns the output amount for the swap.
func (s *V4Swap) AmountOut() *big.Int {
	if s.amount0.Sign() < 0 {
		return new(big.Int).Neg(s.amount0)
	}
	return s.amount1
}

// ParseSwap parses a Uniswap V4 swap log into a V4Swap struct.
// Returns nil if the log is not a valid swap event.
func ParseSwap(log *types.Log) *V4Swap {
	if len(log.Topics) != 3 || len(log.Data) < 64 {
		return nil
	}

	var poolID [32]byte
	copy(poolID[:], log.Topics[1].Bytes())

	amount0 := tools.DecodeSignedInt256(log.Data[:32])
	amount1 := tools.DecodeSignedInt256(log.Data[32:64])

	return &V4Swap{
		poolID:  poolID,
		amount0: amount0,
		amount1: amount1,
	}
}
