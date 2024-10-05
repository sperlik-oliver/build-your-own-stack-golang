package main

import (
	"github.com/gin-gonic/gin"
)

var blockchain Blockchain = Blockchain{
	blocks: []Block{},
}

func main() {
	router := gin.Default()

	router.POST("/health", healthHandler)
	router.GET("/", getBlockchainHandler)
	router.POST("/", createBlockHandler)

	if len(blockchain.blocks) == 0 {
		go blockchain.GenerateOrigin()
	}

	router.Run()

}
