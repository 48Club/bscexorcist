package uniswapv4

import (
	"github.com/48Club/bscexorcist/protocols/tools"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type V4Swap struct {
	poolId  [32]byte // the poolId is a 32-byte identifier for the pool, use it as a unique pool identifier
	amount0 *big.Int
	amount1 *big.Int
}

// Pool use poolId as a unique pool identifier
func (s *V4Swap) Pool() common.Address {
	// 直接使用 poolId 的前 20 字节作为虚拟地址
	// 这样同一个池子的所有 swap 会返回相同的值
	return common.BytesToAddress(s.poolId[:20])
}

func (s *V4Swap) IsToken0To1() bool {
	// amount0 > 0 表示 token0 流入池子（卖出 token0，买入 token1）
	return s.amount0.Sign() > 0
}

func (s *V4Swap) AmountIn() *big.Int {
	if s.amount0.Sign() > 0 {
		return s.amount0
	}
	return new(big.Int).Neg(s.amount1)
}

func (s *V4Swap) AmountOut() *big.Int {
	if s.amount0.Sign() < 0 {
		return new(big.Int).Neg(s.amount0)
	}
	return s.amount1
}

func ParseSwap(log *types.Log) *V4Swap {
	if len(log.Topics) != 3 || len(log.Data) < 64 {
		return nil
	}

	var poolId [32]byte
	copy(poolId[:], log.Topics[1].Bytes())

	amount0 := tools.DecodeSignedInt256(log.Data[:32])
	amount1 := tools.DecodeSignedInt256(log.Data[32:64])

	return &V4Swap{
		poolId:  poolId,
		amount0: amount0,
		amount1: amount1,
	}
}
