# BSC Exorcist 🛡️

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/bsc-exorcist.svg)](https://pkg.go.dev/github.com/yourusername/bsc-exorcist)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Lightweight Go SDK for detecting sandwich attacks in BSC transaction bundles

## 🎯 Purpose

BSC Exorcist provides a simple and efficient way to detect sandwich attack patterns in transaction bundles. This SDK is
**required** for builders submitting bids to 48Club validators, ensuring MEV protection and fair transaction ordering.

Contributions welcome! Help us improve functionality or add support for more protocols.

## ⚡ Quick Start

### Installation

```bash
go get github.com/48Club/bscexorcist
```

### Basic Usage

```go
import "github.com/48Club/bscexorcist"

// Detect sandwich attacks in transaction bundle
err := bscexorcist.DetectSandwichForBundle(transactionsLogs)
if err != nil {
// Sandwich attack detected - handle accordingly
log.Printf("Attack detected: %v", err)
}
```

## 🔍 How It Works

The SDK analyzes DEX swap patterns across transaction bundles to identify sandwich attacks:

```
Bundle Structure:
┌─────────────┐
│   TX 1: Buy │  ← Attacker front-runs
├─────────────┤
│   TX 2: Buy │  ← Victim transaction
├─────────────┤
│  TX 3: Sell │  ← Attacker back-runs
└─────────────┘
```

### Detection Patterns

- **Buy-Buy-Sell**: Front-run and back-run pattern
- **Sell-Sell-Buy**: Reverse sandwich pattern

## 📋 Requirements

- Go 1.21 or higher
- Minimum 3 transactions per bundle
- Transaction logs, separated by each transaction in the bundle

## 🏗️ For Builders

### Integration with 48Club Validators

Builders **MUST** implement sandwich detection before submitting bids:

## 📊 Supported Protocols

| Protocol       | Status      | Event Signatures                |
|----------------|-------------|---------------------------------|
| Uniswap V2     | ✅ Supported | `0xd78ad95f...` `0x606ecd02...` |
| Uniswap V3     | ✅ Supported | `0xc42079f9...` `0x19b47279...` |
| Uniswap V4     | ✅ Supported | `0x40e9cecb...`                 |
| PancakeSwap V2 | ✅ Supported | Compatible                      |
| PancakeSwap V3 | ✅ Supported | Compatible                      |
| PancakeSwap V4 | ✅ Supported | `0x04206ad2...`                 |
| DODOSwap       | ✅ Supported | `0xc2c0245e...`                 |
| FourMeme       | ✅ Supported | `0x7db52723...` `0x0a5575b3...` |

## 🔗 Resources

- [48Club Validator Documentation](https://docs.48.club/48-validators/for-mev-builders)

---

<p align="center">
  Built with ❤️ for the BSC ecosystem
</p>
