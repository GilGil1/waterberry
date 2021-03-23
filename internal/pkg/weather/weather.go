package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Weatherforecast struct {
	ForecastTimeEpoch int `json:"dt"`
	ForecastTime      time.Time
	Main              WeatherforecastMain      `json:"main"`
	WeatherList       []WeatherforecastWeather `json:"weather"`
	Rain              WeatherforecastRain      `json:"rain"`
}

func (wf Weatherforecast) GetTime() time.Time {
	return time.Unix(int64(wf.ForecastTimeEpoch), 0)
}

type WeatherforecastMain struct {
	Temperature float32 `json:"temp"`
	FeelsLike   float32 `json:"feels_like"`
	TempMin     float32 `json:"temp_min"`
	TempMax     float32 `json:"temp_max"`
	// "pressure": 1014,
	// "sea_level": 1014,
	// "grnd_level": 1009,
	// "humidity": 58,
	// "temp_kf": 0.6
}

type WeatherforecastRain struct {
	Probabilityt3h float32 `json:"3h"`
}

type WeatherforecastWeather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type WeatherInfo struct {
	Code                     string `json:"cod"`
	Message                  int    `json:"message"`
	Count                    int    `json:"cnt"`
	NextRainTime             time.Time
	NextRainStartDescription string
	NextRainStopDescription  string
	NextRainStopTime         time.Time
	NextRainDurationMinutes  int
	WeatherForecastList      []Weatherforecast `json:"list"`
}

func NewWeatherInfo(data []byte) (WeatherInfo, error) {
	var weatherInfo WeatherInfo
	err := json.Unmarshal(data, &weatherInfo)
	if err != nil {
		fmt.Printf("An error occured: %v", err)
	}
	// Calculate date in time:
	for i := 0; i < len(weatherInfo.WeatherForecastList); i++ {
		weatherInfo.WeatherForecastList[i].ForecastTime = weatherInfo.WeatherForecastList[i].GetTime()
	}

	/// calcluate when it rains and stops
	for i := 0; i < len(weatherInfo.WeatherForecastList); i++ {
		description := weatherInfo.WeatherForecastList[i].WeatherList[0].Description
		if strings.Contains(description, "rain") {
			weatherInfo.NextRainTime = weatherInfo.WeatherForecastList[i].ForecastTime
			weatherInfo.NextRainStartDescription = description
			for j := i; j < len(weatherInfo.WeatherForecastList); j++ {
				stopDescription := weatherInfo.WeatherForecastList[j].WeatherList[0].Description
				if !strings.Contains(stopDescription, "rain") {
					weatherInfo.NextRainStopTime = weatherInfo.WeatherForecastList[j].ForecastTime
					weatherInfo.NextRainDurationMinutes = int(weatherInfo.NextRainStopTime.Sub(weatherInfo.NextRainTime).Minutes())
					weatherInfo.NextRainStopDescription = weatherInfo.WeatherForecastList[j-1].WeatherList[0].Description
					break
				}
			}
			break
		}
	}
	return weatherInfo, err
}

func (wi WeatherInfo) PrettyPrint() string {
	data, err := json.Marshal(wi)
	if err != nil {
		fmt.Printf("An error occured: %v", err)
	}
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, data, "", "\t")
	if error != nil {
		fmt.Printf("An error occured: %v", err)
	}
	return prettyJSON.String()
}
