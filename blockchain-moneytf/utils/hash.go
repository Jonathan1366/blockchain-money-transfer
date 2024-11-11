package utils

import (
	"crypto/sha256"
	"fmt"
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

	for{
		block.Nonce = nonce
		hash:= GenerateHash(fmt.Sprintf("%d%s%s%d", block.TransactionId, block.PreviousHash, block.Timestamp, block.Nonce))
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

