// Package uniswapv2 provides swap event parsing for Uniswap V2 and compatible protocols.
package uniswapv2

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

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
	delta0 := new(big.Int).Sub(s.amount0Out, s.amount0In) // > 0 means token0 is sent out
	delta1 := new(big.Int).Sub(s.amount1Out, s.amount1In) // > 0 means token1 is sent out
	return delta0.Sign() < 0 && delta1.Sign() > 0
}

// AmountIn returns the input amount for the swap.
func (s *V2Swap) AmountIn() *big.Int {
	delta0 := new(big.Int).Sub(s.amount0Out, s.amount0In)
	delta1 := new(big.Int).Sub(s.amount1Out, s.amount1In)
	if delta0.Sign() < 0 {
		return new(big.Int).Abs(delta0)
	} else {
		return new(big.Int).Abs(delta1)
	}
}

// AmountOut returns the output amount for the swap.
func (s *V2Swap) AmountOut() *big.Int {
	delta0 := new(big.Int).Sub(s.amount0Out, s.amount0In)
	delta1 := new(big.Int).Sub(s.amount1Out, s.amount1In)
	if delta0.Sign() > 0 {
		return delta0
	}
	return delta1
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
