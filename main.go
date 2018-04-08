package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

const (
	Host                = "http://88.198.13.202:9051"
	UrlParallelRequests = 20
)

func main() {
	ergoNodeClient := initializeClient()
	//resp, err := ergoNodeClient.GetBlocks(&GetBlocksRequest{Offset: 10000})
	//resp, err := ergoNodeClient.GetBlock("Vumq5gex8Ty3TuAk8Xxxc9UmRgRd64pnRxvV3PM7Q4Q")
	//if err != nil {
	//	panic(err)
	//}
	//decoded, err := DecodeToBig([]byte(resp.Header.Id))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%b\n", decoded)
	resp, err := ergoNodeClient.GetBlocks(&GetBlocksRequest{Limit: 50, Offset: 0})
	if err != nil {
		panic(err)
	}

	blocks := make([]Block, len(resp.BlocksIds))
	wg := &sync.WaitGroup{}
	semaphore := make(chan interface{}, UrlParallelRequests)

	for i, blockId := range resp.BlocksIds {
		wg.Add(1)
		go func(i int, blockId BlockId) {
			defer wg.Done()
			semaphore <- struct{}{}
			block, err := ergoNodeClient.GetBlock(string(blockId))
			if err != nil {
				panic(err)
			}
			blocks[i] = *block
			<-semaphore
		}(i, blockId)
	}
	wg.Wait()
	proofs, lastBlocks := Prove(&Chain{blocks})
	fmt.Println("================================================================")
	fmt.Println("Proofs:")
	for _, proof := range proofs {
		fmt.Printf("%+v\n", proof.Header)
	}
	fmt.Println("================================================================")
	fmt.Println("LastBlocks:")
	for _, lastBlock := range lastBlocks {
		fmt.Printf("%+v\n", lastBlock.Header)
	}

	//verified, proofs := Verify(proofs, lastBlocks)
}

func initializeClient() *ErgoNodeClient {
	httpClient := http.Client{Timeout: time.Second}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar().With(zap.String("logger", "ErgoNodeClient"))

	client := ErgoNodeClient{
		URL:    Host,
		Logger: logger,
		Client: &httpClient,
	}
	return &client
}
