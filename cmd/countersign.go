package cmd

import (
	"crypto/sha256"
	"encoding/binary"
	"strconv"
	"strings"
	"time"

	"github.com/hwhang0917/countersign/internal/constants"
)

func generateSHA256Hash(text string) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	r := hasher.Sum(nil)
	data := binary.BigEndian.Uint64(r)
	return data
}

func getCurrentInterval() int64 {
	return time.Now().Unix() / int64(constants.INTERVAL)
}

func GenerateOTP(secretKey string, askKey string) string {
	interval := getCurrentInterval()
	seed := strings.Join([]string{secretKey, askKey, strconv.FormatInt(interval, 10)}, ":")
	randomValue := generateSHA256Hash(seed)

	wordListSize := len(constants.WORDS)
	wordIndex := randomValue % uint64(wordListSize)

	return constants.WORDS[wordIndex]
}
