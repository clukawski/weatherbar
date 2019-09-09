package main

import (
	"fmt"
	"os"
)

func main() {
	// Exit on incorrect number of arguments
	if len(os.Args) < 4 {
		os.Exit(1)
	}

	// todo: use flag later
	apiKey := os.Args[1]
	city := os.Args[2]
	countryCode := os.Args[3]

	// Fetch weather data
	wr, err := GetWeather(apiKey, city, countryCode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Display weather data
	fmt.Println(wr.String())
}
