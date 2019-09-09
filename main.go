package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		os.Exit(1)
	}

	// todo: use flag later
	apiKey := os.Args[1]
	city := os.Args[2]
	countryCode := os.Args[3]

	wr, err := GetWeather(apiKey, city, countryCode)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println(wr.String())
}
