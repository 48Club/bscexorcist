// Package uniswapv4 provides swap event parsing for Uniswap V4 and compatible protocols.
package uniswap_v4

import (
	"math/big"

	"github.com/48Club/bscexorcist/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SwapEventSignatures for Uniswap V4
var SwapEventSignatures = map[common.Hash]bool{
	common.HexToHash("0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f"): true, // uniswap-v4
	common.HexToHash("0x04206ad2b7c0f463bff3dd4f33c5735b0f2957a351e4f79763a4fa9e775dd237"): true, // pancake-infinity-cl
}

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

	return &V4Swap{
		poolID:  utils.GetBytes32(log.Topics[1].Bytes()),
		amount0: utils.DecodeSignedInt256(log.Data[:32]),
		amount1: utils.DecodeSignedInt256(log.Data[32:64]),
	}
}
