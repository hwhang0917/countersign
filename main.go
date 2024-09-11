package main

import (
	"flag"
)

func main() {
	nFlag := flag.Int("n", 0, "an int")
	flag.Parse()
	println("nFlag:", *nFlag)
}
