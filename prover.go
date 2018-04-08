package main

const (
	k = 6
	m = 6
)

func IsHaveSuperchainQuality(superchain *Chain, alpha []Block, level int) bool {
	return false
}

func IsHaveMultilevelQuality(superchain *Chain, alpha []Block, level int) bool {
	return false
}

func IsGoodSuperchain(superchain *Chain, alpha []Block, level int) bool {
	return IsHaveSuperchainQuality(superchain, alpha, level) && IsHaveMultilevelQuality(superchain, alpha, level)
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
