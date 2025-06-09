package protocols

import (
	"github.com/48Club/bscexorcist/protocols/uniswapv2"
	"github.com/48Club/bscexorcist/protocols/uniswapv3"
	"github.com/48Club/bscexorcist/protocols/uniswapv4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

// SwapEvent represents a DEX swap event
type SwapEvent interface {
	Pool() common.Address
	IsToken0To1() bool
	AmountIn() *big.Int
	AmountOut() *big.Int
}

var (
	uniswapV2SwapSignatures = map[common.Hash]bool{
		// Swap(address,uint256,uint256,uint256,uint256,address)
		common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"): true,
		// Swap with fee precision
		common.HexToHash("0x606ecd02b3e3b4778f8e97b2e03351de14224efaa5fa64e62200afc9395c2499"): true,
	}

	uniswapV3SwapSignatures = map[common.Hash]bool{
		// Swap(address,address,int256,int256,uint160,uint128,int24)
		common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"): true,
		// Swap with protocol fees
		common.HexToHash("0x19b47279256b2a23a1665c810c8d55a1758940ee09377d4f8d26497a3577dc83"): true,
	}

	uniswapV4SwapSignature = common.HexToHash("0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f")
)

// ParseSwapEvents extracts swap events from transaction logs
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
		} else if uniswapV4SwapSignature == signature {
			swap = uniswapv4.ParseSwap(log)
		}

		if swap != nil {
			swaps = append(swaps, swap)
		}
	}

	return swaps
}
