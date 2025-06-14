// Package bscexorcist provides sandwich attack detection for BSC transaction bundles.
package bscexorcist

import (
	"fmt"
	"github.com/48Club/bscexorcist/protocols"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// DetectSandwichForBundle analyzes a bundle of transaction logs to identify potential sandwich attacks.
// Returns an error if a sandwich pattern is detected in any pool within the bundle.
func DetectSandwichForBundle(bundleLogs [][]*types.Log) error {
	if len(bundleLogs) < 3 {
		return nil
	}

	poolDirections := make(map[common.Address][]bool)
	for _, txLogs := range bundleLogs {
		for _, swap := range protocols.ParseSwapEvents(txLogs) {
			poolID := swap.PairID()
			poolDirections[poolID] = append(poolDirections[poolID], swap.IsToken0To1())
		}
	}

	for pool, directions := range poolDirections {
		if hasSandwichPattern(directions) {
			return fmt.Errorf("sandwich attack detected on pool: %s", pool.Hex())
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
