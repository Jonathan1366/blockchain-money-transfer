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
	prefix := strings.Repeat("0", difficulty)
	for !strings.HasPrefix(block.Hash, prefix) {
		block.Nonce++
		block.Hash = GenerateHash(fmt.Sprintf("%d%s%s%d", block.TransactionId, block.PreviousHash, block.Timestamp, block.Nonce))
	}
}

func GenerateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

