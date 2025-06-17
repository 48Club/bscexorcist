// Package bscexorcist provides sandwich attack detection for BSC transaction bundles.
package bscexorcist

import (
	"fmt"
	"github.com/48Club/bscexorcist/protocols"
	"github.com/48Club/bscexorcist/protocols/liquiditychange"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type SwapDirectionOrType int

const (
	SwapFrom0To1 SwapDirectionOrType = 1
	SwapFrom1To0 SwapDirectionOrType = 0
	LiqChange    SwapDirectionOrType = 2
)

// DetectSandwichForBundle analyzes a bundle of transaction logs to identify potential sandwich attacks.
// Returns an error if a sandwich pattern is detected in any pool within the bundle.
func DetectSandwichForBundle(bundleLogs [][]*types.Log) error {
	if len(bundleLogs) < 3 {
		return nil
	}

	poolDirections := make(map[common.Address][]bool)

	poolDirectionsWithLiqType := make(map[common.Address][]SwapDirectionOrType)
	includeMint := false

	for _, txLogs := range bundleLogs {
		for _, swap := range protocols.ParseSwapEvents(txLogs) {
			poolID := swap.PairID()
			// poolDirections[poolID] = append(poolDirections[poolID], swap.IsToken0To1())

			if _, ok := swap.(*liquiditychange.LiquidityChange); ok {
				includeMint = true
				poolDirectionsWithLiqType[poolID] = append(poolDirectionsWithLiqType[poolID], LiqChange)
			} else {
				poolDirections[poolID] = append(poolDirections[poolID], swap.IsToken0To1())

				if swap.IsToken0To1() {
					poolDirectionsWithLiqType[poolID] = append(poolDirectionsWithLiqType[poolID], SwapFrom0To1)
				} else {
					poolDirectionsWithLiqType[poolID] = append(poolDirectionsWithLiqType[poolID], SwapFrom1To0)
				}

			}
		}
	}

	for pool, directions := range poolDirections {
		if hasSandwichPattern(directions) {
			return fmt.Errorf("sandwich attack detected on pool: %s", pool.Hex())
		}
	}

	if includeMint {
		for pool, directions := range poolDirectionsWithLiqType {
			if hasLiquiditySandwichPattern(directions) {
				return fmt.Errorf("sandwich attack detected on pool: %s", pool.Hex())
			}
		}
	}

	return nil
}

// hasSandwichPattern checks if swap directions form a sandwich attack pattern.
func hasSandwichPattern(directions []bool) bool {
	n := len(directions)
	if n < 3 {
		return false
	}

	// Look for Buy-Buy-Sell or Sell-Sell-Buy patterns
	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				// Buy-Buy-Sell pattern
				if directions[i] && directions[j] && !directions[k] {
					return true
				}
				// Sell-Sell-Buy pattern
				if !directions[i] && !directions[j] && directions[k] {
					return true
				}
			}
		}
	}

	return false
}

func hasLiquiditySandwichPattern(directions []SwapDirectionOrType) bool {
	n := len(directions)
	if n < 3 {
		return false
	}

	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				// Buy-Mint/Burn-Sell pattern or Sell-Mint/Burn-Buy pattern
				if directions[i] == SwapFrom0To1 && directions[j] == LiqChange && directions[k] == SwapFrom1To0 {
					return true
				}
				if directions[i] == SwapFrom1To0 && directions[j] == LiqChange && directions[k] == SwapFrom0To1 {
					return true
				}
			}
		}
	}

	return false
}
