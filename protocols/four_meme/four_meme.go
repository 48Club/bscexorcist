// Package four_meme provides swap event parsing for fourmeme protocols.
package four_meme

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// SwapEventSignatures for FourMemeSwap
	SwapEventSignatures = map[common.Hash]bool{
		swapBuySignature:  true,
		swapSellSignature: true,
	}

	swapBuySignature  = common.HexToHash("0x7db52723a3b2cdd6164364b3b766e65e540d7be48ffa89582956d8eaebe62942")
	swapSellSignature = common.HexToHash("0x0a5575b3648bae2210cee56bf33254cc1ddfbc7bf637c0af2ac18b14fb1bae19")
)

// FourMemeSwap implements SwapEvent for FourMemeSwap protocol.
type FourMemeSwap struct {
	tokenID common.Address
	buySide bool
}

// PairID returns a pseudo-address derived from the first 10 bytes of each token in the pair.
func (s *FourMemeSwap) PairID() common.Address {
	return s.tokenID
}

// IsToken0To1 returns true if the swap direction is token0 -> token1.
func (s *FourMemeSwap) IsToken0To1() bool {
	return s.buySide
}

// AmountIn returns the input amount for the swap.
func (s *FourMemeSwap) AmountIn() *big.Int {
	return big.NewInt(0)
}

// AmountOut returns the output amount for the swap.
func (s *FourMemeSwap) AmountOut() *big.Int {
	return big.NewInt(0)
}

// ParseSwap parses a FourmemeSwap log into a FourmemeSwap struct.
// Returns nil if the log is not a valid swap event.
func ParseSwap(log *types.Log) *FourMemeSwap {
	if len(log.Topics) != 1 || len(log.Data) < 32 {
		return nil
	}
	return &FourMemeSwap{
		tokenID: common.BytesToAddress(log.Data[:32]),
		buySide: log.Topics[0] == swapBuySignature,
	}
}
