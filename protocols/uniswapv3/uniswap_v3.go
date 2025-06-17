// Package uniswapv3 provides swap event parsing for Uniswap V3 and compatible protocols.
package uniswapv3

import (
	"github.com/48Club/bscexorcist/protocols/tools"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// SwapEventSignatures for Uniswap V3
var SwapEventSignatures = map[common.Hash]bool{
	common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"): true,
	common.HexToHash("0x19b47279256b2a23a1665c810c8d55a1758940ee09377d4f8d26497a3577dc83"): true,
}

// V3Swap implements SwapEvent for Uniswap V3-style pools.
type V3Swap struct {
	pool       common.Address
	amount0    *big.Int
	amount1    *big.Int
	zeroForOne bool
}

// PairID returns the pool address.
func (s *V3Swap) PairID() common.Address {
	return s.pool
}

// IsToken0To1 returns true if the swap direction is token0 -> token1.
func (s *V3Swap) IsToken0To1() bool {
	return s.zeroForOne
}

// AmountIn returns the input amount for the swap.
func (s *V3Swap) AmountIn() *big.Int {
	if s.zeroForOne {
		return new(big.Int).Abs(s.amount0)
	}
	return new(big.Int).Abs(s.amount1)
}

// AmountOut returns the output amount for the swap.
func (s *V3Swap) AmountOut() *big.Int {
	if s.zeroForOne {
		return new(big.Int).Abs(s.amount1)
	}
	return new(big.Int).Abs(s.amount0)
}

// ParseSwap parses a Uniswap V3 swap log into a V3Swap struct.
// Returns nil if the log is not a valid swap event.
func ParseSwap(log *types.Log) *V3Swap {
	if len(log.Data) < 160 {
		return nil
	}

	amount0 := tools.DecodeSignedInt256(log.Data[:32])
	amount1 := tools.DecodeSignedInt256(log.Data[32:64])

	return &V3Swap{
		pool:       log.Address,
		amount0:    amount0,
		amount1:    amount1,
		zeroForOne: amount0.Cmp(amount1) > 0,
	}
}
