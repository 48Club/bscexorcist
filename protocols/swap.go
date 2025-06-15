// Package protocols provides unified swap event parsing for supported DEX protocols.
package protocols

import (
	"github.com/48Club/bscexorcist/protocols/dodoswap"
	"github.com/48Club/bscexorcist/protocols/fourmeme"
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

var (
	// Uniswap V2 and compatible swap event signatures
	uniswapV2SwapSignatures = map[common.Hash]bool{
		common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"): true,
		common.HexToHash("0x606ecd02b3e3b4778f8e97b2e03351de14224efaa5fa64e62200afc9395c2499"): true,
	}

	// Uniswap V3 and compatible swap event signatures
	uniswapV3SwapSignatures = map[common.Hash]bool{
		common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"): true,
		common.HexToHash("0x19b47279256b2a23a1665c810c8d55a1758940ee09377d4f8d26497a3577dc83"): true,
	}

	// Uniswap V4 and compatible swap event signature
	uniswapV4SwapSignature = common.HexToHash("0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f")

	// DODOSwap signature for swap events
	dodoSwapSignature = common.HexToHash("0xc2c0245e056d5fb095f04cd6373bc770802ebd1e6c918eb78fdef843cdb37b0f")

	fourMemeSwapSignatures = map[common.Hash]bool{
		common.HexToHash("0x7db52723a3b2cdd6164364b3b766e65e540d7be48ffa89582956d8eaebe62942"): true,
		common.HexToHash("0x0a5575b3648bae2210cee56bf33254cc1ddfbc7bf637c0af2ac18b14fb1bae19"): true,
	}
)

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
		if uniswapV2SwapSignatures[signature] {
			swap = uniswapv2.ParseSwap(log)
		} else if uniswapV3SwapSignatures[signature] {
			swap = uniswapv3.ParseSwap(log)
		} else if signature == uniswapV4SwapSignature {
			swap = uniswapv4.ParseSwap(log)
		} else if signature == dodoSwapSignature {
			swap = dodoswap.ParseSwap(log)
		} else if fourMemeSwapSignatures[signature] {
			swap = fourmeme.ParseSwap(log)
		}

		if swap != nil {
			swaps = append(swaps, swap)
		}
	}

	return swaps
}
