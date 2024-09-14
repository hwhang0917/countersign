package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hwhang0917/countersign/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		fmt.Println("SECRET_KEY is not set")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide ask key")
		fmt.Println("Usage: countersign <ask_key>")
		os.Exit(1)
	}
	askKey := os.Args[1]

	fmt.Print("\033[s")
	for {
		fmt.Print("\033[u\033[K")
		otp, timeLeft := cmd.GenerateOTP(secretKey, askKey)
		fmt.Printf("OTP: %s\nTime left: %d seconds", otp, timeLeft)
		time.Sleep(1 * time.Second)
	}
}
