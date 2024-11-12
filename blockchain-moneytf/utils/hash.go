package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Jonathan1366/blockchain-money-transfer/models"
)

func SignTransaction(privateKey, data string) string {
	hash := sha256.Sum256([]byte(data + privateKey))
	return fmt.Sprintf("%x", hash)
}

// Simple proof-of-work mining function
func MineBlock(block *models.Block, difficulty int) {
	target := strings.Repeat("0", difficulty)
	nonce := 0

	transactionsJSON, err:= json.Marshal(block.Transactions)
	if err != nil {
		log.Printf("Error marshalling transactions: %v", err)
		return
	}

	for{
		block.Nonce = nonce
		hash:= GenerateHash(fmt.Sprintf("%s%s%s%d", transactionsJSON, block.PreviousHash, block.Timestamp, block.Nonce))
		if strings.HasPrefix(hash, target){
			block.Hash=hash
			break 
		}
		nonce++
	}
}

func GenerateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

