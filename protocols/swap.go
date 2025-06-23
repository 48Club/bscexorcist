// Package protocols provides unified swap event parsing for supported DEX protocols.
package protocols

import (
	"math/big"

	"github.com/48Club/bscexorcist/protocols/dodo_swap"
	"github.com/48Club/bscexorcist/protocols/four_meme"
	"github.com/48Club/bscexorcist/protocols/liquidity_change"
	"github.com/48Club/bscexorcist/protocols/uniswap_v2"
	"github.com/48Club/bscexorcist/protocols/uniswap_v3"
	"github.com/48Club/bscexorcist/protocols/uniswap_v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SwapEvent represents a DEX swap event with a unified interface for all supported protocols.
type SwapEvent interface {
	PairID() common.Address
	IsToken0To1() bool
	AmountIn() *big.Int
	AmountOut() *big.Int
}

// ParseSwapEvents extracts swap events from a slice of logs for a single transaction.
// Returns a slice of SwapEvent for all recognized swap events in the logs.
func ParseSwapEvents(logs []*types.Log) []SwapEvent {
	var swaps []SwapEvent

	for _, log := range logs {
		if len(log.Topics) == 0 {
			continue
		}

		var swap SwapEvent
		signature := log.Topics[0]

		if uniswap_v2.SwapEventSignatures[signature] {
			swap = uniswap_v2.ParseSwap(log)
		} else if uniswap_v3.SwapEventSignatures[signature] {
			swap = uniswap_v3.ParseSwap(log)
		} else if uniswap_v4.SwapEventSignatures[signature] {
			swap = uniswap_v4.ParseSwap(log)
		} else if signature == dodo_swap.SwapEventSignature {
			swap = dodo_swap.ParseSwap(log)
		} else if four_meme.SwapEventSignatures[signature] {
			swap = four_meme.ParseSwap(log)
		} else if liquidity_change.LiquidityChangeSignature[signature] {
			swap = liquidity_change.ParseSwap(log)
		}

		if swap != nil {
			swaps = append(swaps, swap)
		}
	}

	return swaps
}
