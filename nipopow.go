package main

const genesisId = "5CBFsnFA67379N8Em59skSQUwDuK3VHoeCZ92DPWaPd7"

type Chain struct {
	blocks []Block
}

func IsValidInterlink(blocks []Block) bool {
	if blocks[0].Header.Interlinks[0] != genesisId {
		// Первый линк всегда генезис
		return false
	}

	for i := len(blocks) - 1; i > 1 ; i-- {
		blockInterlinks := blocks[i].Header.Interlinks
		if blockInterlinks[len(blockInterlinks) - 1] != blocks[i-1].Header.Id {
			return false
		}
	}

	return true
}

func IsValidChain(proofs []Block, lastBlocks []Block) bool {
	if !IsValidInterlink(proofs) {
		return false
	}

	// TODO: валидация интерлинков для последних блоков?
	//if !IsValidInterlink(lastBlocks) {
	//	return false
	//}

	return true
}

// точка входа
func Verify(proofs []Block, lastBlocks []Block) (bool, []Block) {
	ergoNodeClient := initializeClient()
	genesisBlock, err := ergoNodeClient.GetBlock(genesisId)
	if err != nil {
		panic(err)
	}

	proofsChecked := make([]Block, 0)
	proofsChecked = append(proofsChecked, *genesisBlock)

	if IsValidChain(proofs, lastBlocks) {
		proofsChecked = proofs
		return true, proofsChecked
	}

	return false, proofsChecked
}
