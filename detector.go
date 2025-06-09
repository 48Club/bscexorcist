package bscexorcist

import (
	"fmt"
	"github.com/48Club/bscexorcist/protocols"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// DetectSandwichForBundle analyzes transaction logs to identify potential sandwich attacks
func DetectSandwichForBundle(txsLogs [][]*types.Log) error {
	if len(txsLogs) < 3 {
		return nil
	}

	poolSwapDirections := make(map[common.Address][]bool)
	for _, txLogs := range txsLogs {
		for _, swap := range protocols.ParseSwapEvents(txLogs) {
			poolSwapDirections[swap.Pool()] = append(poolSwapDirections[swap.Pool()], swap.IsToken0To1())
		}
	}

	for pool, directions := range poolSwapDirections {
		if containsSandwichPattern(directions) {
			return fmt.Errorf("sandwich attack detected on pool: %s", pool.Hex())
		}
	}

	return nil
}

// containsSandwichPattern checks if swap directions form a sandwich attack pattern
func containsSandwichPattern(directions []bool) bool {
	n := len(directions)
	if n < 3 {
		return false
	}

	// Find pattern: Buy-Buy-Sell or Sell-Sell-Buy
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
