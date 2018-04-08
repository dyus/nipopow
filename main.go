package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	HOST = "http://88.198.13.202:9051"
)

func main() {
	ergoNodeClient := initializeClient()
	//resp, err := ergoNodeClient.GetBlocks(&GetBlocksRequest{Offset: 10000})
	resp, err := ergoNodeClient.GetBlock("Vumq5gex8Ty3TuAk8Xxxc9UmRgRd64pnRxvV3PM7Q4Q")
	if err != nil {
		panic(err)
	}
	decoded, err := DecodeToBig([]byte(resp.Header.Id))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%b\n", decoded)
}

func initializeClient() *ErgoNodeClient {
	httpClient := http.Client{Timeout: time.Second}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar().With(zap.String("logger", "ErgoNodeClient"))

	client := ErgoNodeClient{
		URL:    HOST,
		Logger: logger,
		Client: &httpClient,
	}
	return &client
}
