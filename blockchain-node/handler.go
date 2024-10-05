package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BlockMessage struct {
	Value int `json:"value"`
}

func healthHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"time": time.Now().String(),
	})
}

func createBlockHandler(context *gin.Context) {
	var body BlockMessage
	err := context.BindJSON(&body)
	if err != nil {
		log.Fatalf("Failed parsing body: [%+v]", err)
	}
	blockchain.GenerateBlock(body.Value)
}

func getBlockchainHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"blocks": blockchain.blocks,
	})
}
