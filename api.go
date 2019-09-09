package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

// Various API and formula constants
const (
	ApiURL                = "https://api.openweathermap.org/data/2.5/weather?q=%s,%s&APPID=%s"
	Sonntag90CoefficientA = 17.62
	Sonntag90CoefficientB = 243.12
)

// Output formatting constants
const (
	HumidexFormat   = "Outside: %s, %.1fC, feels like %.1fC"
	WindchillFormat = "Outside: %s, %.1fC, feels like %.1fC"
	DefaultFormat   = "Outside: %s, %.1fC"
)

// Response structure of openweather API
type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"` // m/s, must be converted to km/h to calculate windchill
		Deg   float64 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

// Fetch weather data from openweather and return a pointer to
// the unmarshaled WeatherResponse structure
func GetWeather(apiKey, city, countryCode string) (*WeatherResponse, error) {
	wr := new(WeatherResponse)

	url := fmt.Sprintf(ApiURL, city, countryCode, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBody, wr)
	if err != nil {
		return nil, err
	}

	return wr, nil
}

// Convert kelvin units to celcius
func Celcius(kelvin float64) float64 {
	celcius := kelvin - 273.16
	return celcius
}

// Calculates dewpoint with air temp in celcius and relative humidity
func Dewpoint(temp float64, humidity float64) float64 {
	alpha := func(temp float64, rh float64) float64 {
		return math.Log(rh/100) + ((Sonntag90CoefficientA * temp) / (Sonntag90CoefficientB + temp))
	}

	dewPoint := (Sonntag90CoefficientB * alpha(temp, humidity)) / (Sonntag90CoefficientA - alpha(temp, humidity))

	return dewPoint
}

// Calculates windchill from air temp in celcius and wind speed in km/h
func Windchill(temp, speed float64) float64 {
	windChill := 13.12 + 0.6215*temp - 11.37*math.Pow(speed, 0.16) + 0.3965*temp*math.Pow(speed, 0.16)
	return windChill
}

// Return a string containing basic weather information
func (wr *WeatherResponse) String() string {
	// Do various conversionsa
	windSpeedKph := 3.6 * wr.Wind.Speed
	tempCelcius := Celcius(wr.Main.Temp)

	// Calculate windchill using air temp in celcius and wind speed in km/h
	windChill := Windchill(tempCelcius, windSpeedKph)
	// Calculate dewpoint using Sonntag90 coefficients
	dewPointCelcius := Dewpoint(tempCelcius, float64(wr.Main.Humidity))
	// Convert dewpoint back to kelvin for humidex calculation
	dewPoint := dewPointCelcius + 273.15

	// Calculate humidex using dew point and air temp
	// Todo: document various equation constants
	// Remember 5/9 is zero because 5 and 9 are integer types (do 5.0/9.0 instead)
	e := 6.11 * math.Exp(5417.7530*((1/273.15)-(1/dewPoint)))
	h := (5.0 / 9.0) * (e - 10)
	humidex := wr.Main.Temp + h

	// Determine string output based on air temperature conditions
	var output string
	switch true {
	case tempCelcius < 0:
		output = fmt.Sprintf(WindchillFormat, wr.Weather[0].Main, Celcius(wr.Main.Temp), windChill)
	case tempCelcius > 19:
		output = fmt.Sprintf(HumidexFormat, wr.Weather[0].Main, Celcius(wr.Main.Temp), Celcius(humidex))
	default:
		output = fmt.Sprintf(DefaultFormat, wr.Weather[0].Main, Celcius(wr.Main.Temp))
	}

	return output
}
