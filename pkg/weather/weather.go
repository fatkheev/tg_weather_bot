package weather

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

const weatherAPIURL = "http://api.openweathermap.org/data/2.5/weather"

type WeatherResponse struct {
    Main struct {
        Temp     float64 `json:"temp"`
        Humidity int     `json:"humidity"`
    } `json:"main"`
    Wind struct {
        Speed float64 `json:"speed"`
    } `json:"wind"`
    Name string `json:"name"`
}

func GetWeather(city string) (*WeatherResponse, error) {
    apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("OPENWEATHERMAP_API_KEY не указан в файле .env")
    }

    url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", weatherAPIURL, city, apiKey)
    log.Printf("Запрос данных о погоде с URL: %s", url)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var weatherResponse WeatherResponse
    err = json.Unmarshal(body, &weatherResponse)
    if err != nil {
        return nil, err
    }

    return &weatherResponse, nil
}

func GetWeatherByCoordinates(lat, lon float64) (*WeatherResponse, error) {
    apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("OPENWEATHERMAP_API_KEY не указан в файле .env")
    }

    url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", weatherAPIURL, lat, lon, apiKey)
    log.Printf("Запрос данных о погоде с URL: %s", url)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var weatherResponse WeatherResponse
    err = json.Unmarshal(body, &weatherResponse)
    if err != nil {
        return nil, err
    }

    return &weatherResponse, nil
}