// Package uniswapv2 provides swap event parsing for Uniswap V2 and compatible protocols.
package uniswapv2

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// V2Swap implements SwapEvent for Uniswap V2-style pools.
type V2Swap struct {
	pool       common.Address
	amount0In  *big.Int
	amount1In  *big.Int
	amount0Out *big.Int
	amount1Out *big.Int
}

// PairID returns the pool address.
func (s *V2Swap) PairID() common.Address {
	return s.pool
}

// IsToken0To1 returns true if the swap direction is token0 -> token1.
func (s *V2Swap) IsToken0To1() bool {
	return s.amount0In.Cmp(s.amount1In) > 0
}

// AmountIn returns the input amount for the swap.
func (s *V2Swap) AmountIn() *big.Int {
	if s.amount0In.Cmp(s.amount1In) > 0 {
		return s.amount0In
	}
	return s.amount1In
}

// AmountOut returns the output amount for the swap.
func (s *V2Swap) AmountOut() *big.Int {
	if s.amount0Out.Cmp(s.amount1Out) > 0 {
		return s.amount0Out
	}
	return s.amount1Out
}

// ParseSwap parses a Uniswap V2 swap log into a V2Swap struct.
// Returns nil if the log is not a valid swap event.
func ParseSwap(log *types.Log) *V2Swap {
	if len(log.Data) < 128 {
		return nil
	}

	amount0In := new(big.Int).SetBytes(log.Data[:32])
	amount1In := new(big.Int).SetBytes(log.Data[32:64])
	amount0Out := new(big.Int).SetBytes(log.Data[64:96])
	amount1Out := new(big.Int).SetBytes(log.Data[96:128])

	return &V2Swap{
		pool:       log.Address,
		amount0In:  amount0In,
		amount1In:  amount1In,
		amount0Out: amount0Out,
		amount1Out: amount1Out,
	}
}
