package main

import (
	"fmt"
)

type Chain struct {
	blocks []BlockId
}

func Prove(chain Chain) {
	// check [0]
	B := chain.blocks[0]
	//proofs :=
	fmt.Println(B)
}
