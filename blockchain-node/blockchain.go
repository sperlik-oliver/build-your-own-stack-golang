package main

import "time"

type Blockchain struct {
	blocks []Block
}

func (blockchain *Blockchain) Validate() bool {
	allBlocksValid := true
	for _, block := range blockchain.blocks {
		if !block.IsBlockValid() {
			allBlocksValid = false
		}
	}
	return allBlocksValid
}

func (blockchain *Blockchain) FindPrevious() Block {
	previousBlock := blockchain.blocks[0]
	for _, block := range blockchain.blocks {
		if block.Index > previousBlock.Index {
			previousBlock = block
		}
	}
	return previousBlock
}

func (blockchain *Blockchain) GenerateBlock(value int) {
	previous := blockchain.FindPrevious()
	block := Block{
		Index:        previous.Index + 1,
		Timestamp:    time.Now().String(),
		Value:        value,
		PreviousHash: previous.Hash,
	}
	block.Hash = block.CalculateHash()
	blockchain.blocks = append(blockchain.blocks, block)
}

func (blockchain *Blockchain) GenerateOrigin() {
	block := Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Value:        0,
		PreviousHash: "",
	}
	block.Hash = block.CalculateHash()
	blockchain.blocks = append(blockchain.blocks, block)
}
