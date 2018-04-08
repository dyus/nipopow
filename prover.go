package main

import (
	"math"
	"math/big"
)

const (
	k     = 6
	m     = 6
	delta = 0.5
)

func IsLocallyGoodSuperchain(chain []Block, alpha []Block, level int) bool {
	powered := big.NewInt(int64(2))
	powered.Exp(powered, big.NewInt(int64(level)), nil)

	right := &big.Float{}
	right.Quo(big.NewFloat(float64((1-delta)*float64(len(chain)))), (&big.Float{}).SetInt(powered))
	return big.NewFloat(float64(len(alpha))).Cmp(right) > 0
}

func filterByLevel(chain []Block, level int) *[]Block {
	filteredBlocks := make([]Block, 0)

	for _, block := range chain {
		if block.GetLevel() >= level {
			filteredBlocks = append(filteredBlocks, block)
		}
	}
	return &filteredBlocks
}

func IsHaveSuperchainQuality(chain *Chain, alpha []Block, level int) bool {
	for mS := m; mS < len(alpha); mS++ {
		chainSlice := chain.blocks[len(chain.blocks)-mS:]
		checkChain := make([]Block, 0)
		reversedCheckChain := make([]Block, 0)

		for i, block := range chainSlice {
			if block.GetLevel() >= level {
				checkChain = append(checkChain, block)
			}

			reversedBlock := checkChain[len(checkChain)-1-i]
			if reversedBlock.GetLevel() >= level {
				reversedCheckChain = append(reversedCheckChain, reversedBlock)
			}
		}

		if !IsLocallyGoodSuperchain(checkChain, reversedCheckChain, level) {
			return false
		}
	}
	return true
}

func IsHaveMultilevelQuality(chain *Chain, alpha []Block, level int) bool {
	const k1 = m
	for levelS := level; levelS > 0; levelS-- {
		// Any nested chain
		filteredByLevelSCount := len(*filterByLevel(chain.blocks, levelS))
		if filteredByLevelSCount >= k1 {
			if float64(len(*filterByLevel(chain.blocks, level))) < (1-delta)*float64(filteredByLevelSCount)*math.Pow(2, float64(level-levelS)) {

			}
		}
	}
	return true
}

func IsGoodSuperchain(chain *Chain, alpha []Block, level int) bool {
	return IsHaveSuperchainQuality(chain, alpha, level) && IsHaveMultilevelQuality(chain, alpha, level)
}

func Prove(chain *Chain) ([]Block, []Block) {
	block := chain.blocks[0]
	// mu
	initialLevel := len(chain.blocks[len(chain.blocks)-k].Header.Interlinks)
	proofs := make([]Block, 0)
	for level := initialLevel; level >= 0; level-- {
		alpha := make([]Block, 0, len(chain.blocks)-k)
		for i := 0; i < len(chain.blocks)-k; i++ {
			if block.Header.Height <= chain.blocks[i].Header.Height && chain.blocks[i].GetLevel() >= level {
				alpha = append(alpha, chain.blocks[i])
			}
		}
		proofs = append(proofs, alpha...)

		if IsGoodSuperchain(chain, alpha, level) {
			block = alpha[len(alpha)-m]
		}
	}
	lastBlocks := chain.blocks[len(chain.blocks)-k:]

	return proofs, lastBlocks
}
