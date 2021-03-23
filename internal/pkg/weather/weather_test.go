package weather

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestWeatherSerialize(t *testing.T) {
	filename := "weather.json"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	weatherInfo, err := NewWeatherInfo(content)
	if err != nil {
		t.Errorf("%v", err)
	}
	log.Println("Weather Info:", weatherInfo.PrettyPrint())
}
