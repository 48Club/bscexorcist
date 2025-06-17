// Package liquiditychange provides Mint/Burn event parsing.
package liquiditychange

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// LiquidityChange implements SwapEvent for Mint/Burn Event.
type LiquidityChange struct {
	pool common.Address
}

var (
	LiquidityChangeSignature = map[common.Hash]bool{
		uniswapV2MintSignature: true,
		uniswapV3MintSignature: true,
		uniswapV2BurnSignature: true,
		uniswapV3BurnSignature: true,
	}

	// Mint (index_topic_1 address sender, uint256 amount0, uint256 amount1)
	uniswapV2MintSignature = common.HexToHash("0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f")

	// Mint (address sender, index_topic_1 address owner, index_topic_2 int24 tickLower, index_topic_3 int24 tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
	uniswapV3MintSignature = common.HexToHash("0x7a53080ba414158be7ec69b987b5fb7d07dee101fe85488f0853ae16239d0bde")

	// Burn (index_topic_1 address sender, uint256 amount0, uint256 amount1, index_topic_2 address to)
	uniswapV2BurnSignature = common.HexToHash("0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496")

	// Burn (index_topic_1 address owner, index_topic_2 int24 tickLower, index_topic_3 int24 tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
	uniswapV3BurnSignature = common.HexToHash("0x0c396cd989a39f4459b5fa1aed6a9a8dcdbc45908acfd67e028cd568da98982c")
)

// PairID returns the pool address.
func (s *LiquidityChange) PairID() common.Address {
	return s.pool
}

// In Mint/Burn Event, IsToken0To1 don't make sense, always returns false.
func (s *LiquidityChange) IsToken0To1() bool {
	return false
}

// AmountIn returns 0.
func (s *LiquidityChange) AmountIn() *big.Int {
	return big.NewInt(0)
}

// AmountOut returns 0.
func (s *LiquidityChange) AmountOut() *big.Int {
	return big.NewInt(0)
}

// ParseSwap parses a Mint/Burn log into a LiquidityChange struct.
// Returns nil if the log is not a valid swap event.
func ParseSwap(log *types.Log) *LiquidityChange {
	if len(log.Topics) < 2 || len(log.Data) < 64 {
		return nil
	}

	return &LiquidityChange{
		pool: log.Address,
	}
}
