package main

import "math/big"

const (
	k     = 6
	m     = 6
	delta = 0.5
)

func IsLocallyGoodSuperchain(chain *Chain, alpha []Block, level int) bool {
	powered := big.Int{}.Exp(big.NewInt(int64(2)), big.NewInt(int64(level)), nil)
	right := big.Float{}.Quo(big.NewFloat(float64((1-delta)*float64(len(chain.blocks)))), big.Float{}.SetInt(powered))
	return big.NewFloat(float64(len(alpha))).Cmp(right) > 0
}

func IsHaveSuperchainQuality(chain *Chain, alpha []Block, level int) bool {
	return true
}

func IsHaveMultilevelQuality(chain *Chain, alpha []Block, level int) bool {
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
			if chain.blocks[i].Header.Height >= block.Header.Height && chain.blocks[i].GetLevel() >= level {
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
