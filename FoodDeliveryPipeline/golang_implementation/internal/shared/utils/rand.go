package utils

import (
	"crypto/rand"
	"fmt"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateID(prefix string, length int) (string, error) {
	b := make([]byte, length)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	id := make([]byte, length)

	for i := range id {
		id[i] = chars[int(b[i])%len(chars)]
	}

	return fmt.Sprintf("%s-%s", prefix, string(id)), nil
}

func GenerateOrderID() (string, error) {
	return generateID("ORD", 10)
}

func GenerateTxnID() (string, error) {
	return generateID("TXN", 10)
}

func GenerateGatewayTxnID() (string, error) {
	return generateID("GW", 10)
}
