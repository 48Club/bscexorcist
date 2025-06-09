package uniswapv2

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type V2Swap struct {
	pool       common.Address
	amount0In  *big.Int
	amount1In  *big.Int
	amount0Out *big.Int
	amount1Out *big.Int
}

func (s *V2Swap) Pool() common.Address {
	return s.pool
}

func (s *V2Swap) IsToken0To1() bool {
	return s.amount0In.Cmp(s.amount1In) > 0
}

func (s *V2Swap) AmountIn() *big.Int {
	if s.amount0In.Cmp(s.amount1In) > 0 {
		return s.amount0In
	}
	return s.amount1In
}

func (s *V2Swap) AmountOut() *big.Int {
	if s.amount0Out.Cmp(s.amount1Out) > 0 {
		return s.amount0Out
	}
	return s.amount1Out
}

func ParseSwap(log *types.Log) *V2Swap {
	if len(log.Data) < 128 {
		return nil
	}

	amount0in := new(big.Int).SetBytes(log.Data[:32])
	amount1in := new(big.Int).SetBytes(log.Data[32:64])
	amount0out := new(big.Int).SetBytes(log.Data[64:96])
	amount1out := new(big.Int).SetBytes(log.Data[96:128])

	return &V2Swap{
		pool:       log.Address,
		amount0In:  amount0in,
		amount1In:  amount1in,
		amount0Out: amount0out,
		amount1Out: amount1out,
	}
}
