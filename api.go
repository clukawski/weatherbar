package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

const (
	ApiURL                = "https://api.openweathermap.org/data/2.5/weather?q=Toronto,CA&APPID=%s"
	HumidexFormat         = "Outside: %s, %.1fC, feels like %.1fC"
	WindchillFormat       = "Outside: %s, %.1fC, feels like %.1fC"
	DefaultFormat         = "Outside: %s, %.1fC"
	Sonntag90CoefficientA = 17.62
	Sonntag90CoefficientB = 243.12
)

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
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
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

func GetWeather(apiKey string) (*WeatherResponse, error) {
	wr := new(WeatherResponse)
	resp, err := http.Get(fmt.Sprintf(ApiURL, apiKey))
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

func Celcius(kelvin float64) float64 {
	celcius := kelvin - 273.16
	return celcius
}

func Dewpoint(temp float64, humidity float64) float64 {
	alpha := func(temp float64, rh float64) float64 {
		return math.Log(rh/100) + ((Sonntag90CoefficientA * temp) / (Sonntag90CoefficientB + temp))
	}

	dewPoint := (Sonntag90CoefficientB * alpha(temp, humidity)) / (Sonntag90CoefficientA - alpha(temp, humidity))

	return dewPoint
}

func Windchill(temp, speed float64) float64 {
	windChill := 13.12 + 0.6215*temp - 11.37*math.Pow(speed, 0.16) + 0.3965*math.Pow(temp, 0.16)
	return windChill
}

func (wr *WeatherResponse) String() string {
	tempCelcius := Celcius(wr.Main.Temp)
	dewPointCelcius := Dewpoint(tempCelcius, float64(wr.Main.Humidity))
	dewPoint := dewPointCelcius + 273.15

	// Remember 5/9 is zero because 5 and 9 are integer types (do 5.0/9.0 instead)
	e := 6.11 * math.Exp(5417.7530*((1/273.15)-(1/dewPoint)))
	h := (5.0 / 9.0) * (e - 10)
	humidex := wr.Main.Temp + h
	windChill := Windchill(tempCelcius, wr.Wind.Speed)

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
