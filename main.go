package main

import (
	"fmt"
	"os"

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
		os.Exit(1)
	}
	askKey := os.Args[1]

	otp := cmd.GenerateOTP(secretKey, askKey)
	fmt.Println("OTP: ", otp)
}
