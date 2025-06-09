package uniswapv3

import (
	"github.com/48Club/bscexorcist/protocols/tools"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type V3Swap struct {
	pool       common.Address
	amount0    *big.Int
	amount1    *big.Int
	zeroForOne bool
}

func (s *V3Swap) Pool() common.Address {
	return s.pool
}

func (s *V3Swap) IsToken0To1() bool {
	return s.zeroForOne
}

func (s *V3Swap) AmountIn() *big.Int {
	if s.zeroForOne {
		return new(big.Int).Abs(s.amount0)
	}
	return new(big.Int).Abs(s.amount1)
}

func (s *V3Swap) AmountOut() *big.Int {
	if s.zeroForOne {
		return new(big.Int).Abs(s.amount1)
	}
	return new(big.Int).Abs(s.amount0)
}

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
