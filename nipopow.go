package main

import "go.uber.org/zap"

const genesisId = "5CBFsnFA67379N8Em59skSQUwDuK3VHoeCZ92DPWaPd7"

var logger *zap.SugaredLogger

func init() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger = zapLogger.Sugar().With(zap.String("logger", "ErgoNodeClient"))
}

type Chain struct {
	blocks []Block
}

func IsValidInterlink(blocks []Block) bool {
	if blocks[0].Header.Interlinks[0] != genesisId {
		// Первый линк всегда генезис
		return false
	}

	for i := len(blocks) - 1; i > 1; i-- {
		blockInterlinks := blocks[i].Header.Interlinks
		if blockInterlinks[len(blockInterlinks)-1] != blocks[i-1].Header.Id {
			return false
		}
		logger.Info("Block %s has valid interlink", blocks[i].Header.Id)
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
	logger.Info("Getting genesis block...")
	ergoNodeClient := initializeClient()
	genesisBlock, err := ergoNodeClient.GetBlock(genesisId)
	if err != nil {
		panic(err)
	}

	proofsChecked := make([]Block, 0)
	proofsChecked = append(proofsChecked, *genesisBlock)

	logger.Info("Start validating proofs and last blocks...")

	if IsValidChain(proofs, lastBlocks) {
		proofsChecked = proofs

		logger.Info("Chain is valid")

		return true, proofsChecked
	}

	return false, proofsChecked
}
