package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	Index        int
	Timestamp    string
	Value        int
	Hash         string
	PreviousHash string
}

func (block *Block) ToString() string {
	return string(block.Index) + block.Timestamp + fmt.Sprint(block.Value) + block.Hash + block.PreviousHash
}

func (block *Block) CalculateHash() string {
	hash := sha256.New()
	hash.Write([]byte(block.ToString()))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (block *Block) IsBlockValid() bool {
	previous := blockchain.FindPrevious()
	return previous.Hash == block.PreviousHash
}
