// Package dodoswap provides swap event parsing for DODOSwap protocols.
package dodoswap

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// DODOSwap implements SwapEvent for DODOSwap protocol.
type DODOSwap struct {
	poolID     common.Address
	tokenFrom  common.Address
	tokenTo    common.Address
	amountFrom *big.Int
	amountTo   *big.Int
}

// PairID returns a pseudo-address derived from the first 10 bytes of each token in the pair.
func (s *DODOSwap) PairID() common.Address {
	return s.poolID
}

// IsToken0To1 returns true if the swap direction is token0 -> token1.
func (s *DODOSwap) IsToken0To1() bool {
	return isTokenAFirst(s.tokenFrom, s.tokenTo)
}

// AmountIn returns the input amount for the swap.
func (s *DODOSwap) AmountIn() *big.Int {
	if isTokenAFirst(s.tokenFrom, s.tokenTo) {
		// tokenFrom is token0, so amountFrom is the input amount
		return new(big.Int).Set(s.amountFrom)
	}
	return new(big.Int).Set(s.amountTo)
}

// AmountOut returns the output amount for the swap.
func (s *DODOSwap) AmountOut() *big.Int {
	if isTokenAFirst(s.tokenFrom, s.tokenTo) {
		// tokenFrom is token0, so amountFrom is the input amount
		return new(big.Int).Set(s.amountTo)
	}
	return new(big.Int).Set(s.amountFrom)
}

// ParseSwap parses a DODOSwap log into a DODOSwap struct.
// Returns nil if the log is not a valid swap event.
func ParseSwap(log *types.Log) *DODOSwap {
	if len(log.Topics) != 1 || len(log.Data) < 128 {
		return nil
	}

	fromToken := common.BytesToAddress(log.Data[:32])
	toToken := common.BytesToAddress(log.Data[32:64])

	return &DODOSwap{
		poolID:     calcPoolID(fromToken, toToken),
		tokenFrom:  fromToken,
		tokenTo:    toToken,
		amountFrom: new(big.Int).SetBytes(log.Data[64:96]),
		amountTo:   new(big.Int).SetBytes(log.Data[96:128]),
	}
}

func isTokenAFirst(tkA, tkB common.Address) bool {
	tokenARep := new(big.Int).SetBytes(tkA.Bytes())
	tokenBRep := new(big.Int).SetBytes(tkB.Bytes())

	return tokenARep.Cmp(tokenBRep) < 0
}

// calcPoolID returns a pseudo-address derived from the first 10 bytes of each sorted tokenFrom/tokenTo.
func calcPoolID(tkA, tkB common.Address) (pool common.Address) {
	tokenARep := new(big.Int).SetBytes(tkA.Bytes())
	tokenBRep := new(big.Int).SetBytes(tkB.Bytes())
	if tokenARep.Cmp(tokenBRep) < 0 {
		copy(pool[:], tkA.Bytes()[:10])
		copy(pool[10:], tkB.Bytes()[:10])
		return pool
	} else {
		copy(pool[:], tkB.Bytes()[:10])
		copy(pool[10:], tkA.Bytes()[:10])
	}
	return pool
}
