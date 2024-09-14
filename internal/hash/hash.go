package hash

import (
	"crypto/sha256"
	"encoding/binary"
	"strconv"
	"strings"
	"time"

	"github.com/hwhang0917/countersign/internal/repository"
	"github.com/hwhang0917/countersign/internal/services"
	"github.com/hwhang0917/countersign/pkg/config"
	"github.com/hwhang0917/countersign/pkg/database"
)

func generateSHA256Hash(text string) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	r := hasher.Sum(nil)
	data := binary.BigEndian.Uint64(r)
	return data
}

func getTimestamp() int64 {
	return time.Now().Unix()
}

func getCurrentInerval(interval int64) int64 {
	timestamp := getTimestamp()
	currentInverval := timestamp / interval
	return currentInverval
}

func GetTimeLeft(interval int64) int64 {
	timestamp := getTimestamp()
	timeLeft := interval - (timestamp % interval)
	return timeLeft
}

func GetOTP(askKey string) string {
	interval := config.GetInterval()
	currentInterval := getCurrentInerval(int64(interval))
	seed := strings.Join([]string{askKey, strconv.FormatInt(currentInterval, 10)}, ":")
	hash := generateSHA256Hash(seed)

	wordRepository := repository.NewWordRepository(database.GetDB())
	wordService := services.NewWordService(wordRepository)
	wordListCount, err := wordService.GetWordListCount()
	if err != nil {
		panic(err)
	}
	if wordListCount == 0 {
		return "no word"
	}
	wordId := int(hash) % int(wordListCount)

	word, err := wordService.GetWordByID(uint(wordId))
	if err != nil {
		panic(err)
	}
	return word
}
