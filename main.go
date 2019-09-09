package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	apiKey := os.Args[1]
	wr, err := GetWeather(apiKey)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(wr.String())
}
