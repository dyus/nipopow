package main

import (
	"go.uber.org/zap"
	"fmt"
)

const (
	HOST = "http://88.198.13.202:9051"
)

func main() {
	client := initializeClient()
	resp, err := client.GetBlocks(&GetBlocksRequest{Offset: 10000})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", resp)
}

func initializeClient() *ErgoNodeClient {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar().With(zap.String("logger", "ErgoNodeClient"))
	client := ErgoNodeClient{URL: HOST, Logger: logger}
	return &client
}
