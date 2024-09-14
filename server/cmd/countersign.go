package cmd

import (
	"crypto/sha256"
	"encoding/binary"
	"strconv"
	"strings"
	"time"

	"github.com/hwhang0917/countersign/internal/constants"
	"github.com/hwhang0917/countersign/internal/words"
)

func generateSHA256Hash(text string) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	r := hasher.Sum(nil)
	data := binary.BigEndian.Uint64(r)
	return data
}

func getCurrentInterval() (int64, int64) {
	timeLeft := int64(constants.INTERVAL) - time.Now().Unix()%int64(constants.INTERVAL)
	interval := time.Now().Unix() / int64(constants.INTERVAL)
	return interval, timeLeft
}

func GenerateOTP(secretKey string, askKey string) (string, int64) {
	interval, timeLeft := getCurrentInterval()
	seed := strings.Join([]string{secretKey, askKey, strconv.FormatInt(interval, 10)}, ":")
	randomValue := generateSHA256Hash(seed)

	wordListSize := len(words.WORD_LIST)
	wordIndex := randomValue % uint64(wordListSize)

	return words.WORD_LIST[wordIndex], timeLeft
}
