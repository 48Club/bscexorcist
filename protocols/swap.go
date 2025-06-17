// Package protocols provides unified swap event parsing for supported DEX protocols.
package protocols

import (
	"github.com/48Club/bscexorcist/protocols/dodoswap"
	"github.com/48Club/bscexorcist/protocols/fourmeme"
	"github.com/48Club/bscexorcist/protocols/liquiditychange"
	"github.com/48Club/bscexorcist/protocols/uniswapv2"
	"github.com/48Club/bscexorcist/protocols/uniswapv3"
	"github.com/48Club/bscexorcist/protocols/uniswapv4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
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

		signature := log.Topics[0]

		var swap SwapEvent
		if uniswapv2.SwapEventSignatures[signature] {
			swap = uniswapv2.ParseSwap(log)
		} else if uniswapv3.SwapEventSignatures[signature] {
			swap = uniswapv3.ParseSwap(log)
		} else if signature == uniswapv4.SwapEventSignature {
			swap = uniswapv4.ParseSwap(log)
		} else if signature == dodoswap.SwapEventSignature {
			swap = dodoswap.ParseSwap(log)
		} else if fourmeme.SwapEventSignatures[signature] {
			swap = fourmeme.ParseSwap(log)
		} else if liquiditychange.LiquidityChangeSignature[signature] {
			swap = liquiditychange.ParseSwap(log)
		}

		if swap != nil {
			swaps = append(swaps, swap)
		}
	}

	return swaps
}
