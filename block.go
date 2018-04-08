package main

type BlockId string

type BlockHeader struct {
	Difficulty        string    `json:"difficulty"`
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
