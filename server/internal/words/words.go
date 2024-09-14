package words

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var WORD_LIST []string

func Init() {
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		WORD_LIST = append(WORD_LIST, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(WORD_LIST) < 10 {
		fmt.Println("Word list must have at least 10 words")
		os.Exit(1)
	}
}
