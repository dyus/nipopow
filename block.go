package main

import (
	"math/big"
)

type BlockId string

type BlockHeader struct {
	Difficulty        int       `json:"difficulty,string"`
	Votes             string    `json:"votes"`
	Timestamp         int       `json:"timestamp"`
	Nonce             int       `json:"nonce"`
	StateRoot         string    `json:"stateRoot"`
	Height            int       `json:"height"`
	NBits             int       `json:"nBits"`
	Id                BlockId   `json:"id"`
	Interlinks        []BlockId `json:"interlinks"`
	AdProofsRoot      string    `json:"adProofsRoot"`
	TransactionsRoot  string    `json:"transactionsRoot"`
	EquihashSolutions string    `json:"equihashSolutions"`
	ParentId          BlockId   `json:"parentId"`
}

type BlockADProofs struct {
	HeaderId   string `json:"headerId"`
	ProofBytes string `json:"proofBytes"`
	Digest     string `json:"digest"`
}

type Block struct {
	Header BlockHeader `json:"header"`
	// server typo: adPoofs
	ADProofs BlockADProofs `json:"adPoofs"`
}

func (b *Block) GetLevel() int {
	T := b.Header.Difficulty
	// TODO change type after read from server
	idInt, err := DecodeToBig([]byte(b.Header.Id))
	if err != nil {
		panic(err)
	}
	// TODO normal big math division Log2?
	level := 0
	id := new(big.Float).SetInt(idInt)
	TFloat := big.NewFloat(float64(T))
	for id.Cmp(TFloat) <= 0 {
		level++
		id = new(big.Float).Mul(id, big.NewFloat(float64(2)))
	}
	if level > 0 {
		level--
	}
	return level
}
