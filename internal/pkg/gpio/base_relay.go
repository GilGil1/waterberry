package gpio

import (
	"fmt"
	"time"
	config "waterberry/internal/pkg/config"
)

type IRelay interface {
	SetConfig(id int, name string, pin uint8, timings []config.OpenTimeConfig)
	Init()
	GetPin() uint8
	GetId() int
	GetName() string
	GetSecondsOff() string
	GetCurrentMode() string
	SetOn() error
	SetOff() error
	SetMode(mode string) error
	SetFields(mode string, updateTime time.Time) error
	GetPropertiesMap() map[string]interface{}
}

type BaseRelay struct {
	ID            int
	Pin           uint8
	Name          string
	CurrentStatus string
	UpdateTime    time.Time
	StopTime      time.Time
	timings       []config.OpenTimeConfig
}

func (br *BaseRelay) SetFields(mode string, updateTime time.Time) error {
	if mode == RelayOff || mode == RelayOn {
		br.CurrentStatus = mode
		br.UpdateTime = updateTime
	} else {
		return fmt.Errorf("Mode %s unsupported\n", mode)
	}
	return nil
}
func (br *BaseRelay) GetPropertiesMap() map[string]interface{} {
	var singleMap = make(map[string]interface{})
	singleMap["id"] = br.GetId()
	singleMap["name"] = br.GetName()
	singleMap["pin"] = br.GetPin()
	singleMap["mode"] = br.GetCurrentMode()
	singleMap["seconds_off"] = br.GetSecondsOff()
	return singleMap

}

func (br *BaseRelay) Init() {
	br.startTimers()
}

func (br *BaseRelay) SetConfig(id int, name string, pin uint8, timings []config.OpenTimeConfig) {
	br.ID = id
	br.Name = name
	br.Pin = pin
	br.timings = timings
	if len(timings) > 0 {
		go br.startTimers()
	}
}

func (br *BaseRelay) GetPin() uint8 {
	return br.Pin
}

func (br *BaseRelay) GetId() int {
	return br.ID
}

func (br *BaseRelay) GetName() string {
	return br.Name
}

func (br *BaseRelay) GetSecondsOff() string {
	if br.StopTime.IsZero() {
		return "No Timer Set"
	} else {
		return fmt.Sprintf("%4.0f", time.Until(br.StopTime).Seconds())

	}
}

func (br *BaseRelay) GetCurrentMode() string {
	return br.CurrentStatus
}

func (br *BaseRelay) SetOn() error {

	return br.SetMode(RelayOn)
}

func (br *BaseRelay) SetOff() error {
	return br.SetMode(RelayOff)

}

func (br *BaseRelay) SetMode(mode string) error {
	if mode == RelayOn {
		br.CurrentStatus = mode
		go br.setOnAndTimer(1200)
		printModeChange(br)
	} else if mode == RelayOff {
		br.CurrentStatus = mode
		br.cancelTimer()
		printModeChange(br)
	} else {
		fmt.Errorf("Mode %s unsupportee\n", mode)
	}
	return nil
}
func (br *BaseRelay) cancelTimer() {
	br.StopTime = time.Time{}
}
func (br *BaseRelay) setOnAndTimer(stopinSeconds int) {
	defer func() {
		br.SetMode(RelayOff)
	}()
	br.CurrentStatus = RelayOn
	br.StopTime = time.Now().Add(time.Duration(stopinSeconds) * time.Second)
	timer := time.NewTimer(time.Duration(stopinSeconds) * time.Second)
	<-timer.C
	fmt.Printf("Timer %d fired\n", br.GetId())
}

func (br *BaseRelay) startTimers() {
	ticker := time.NewTicker(time.Millisecond * 5000)
	for t := range ticker.C {
		evalutationTime := time.Now()
		evaluationWeekDay := evalutationTime.Weekday()
		evaluationHour := evalutationTime.Hour()
		evaluationMinute := evalutationTime.Minute()
		fmt.Println("Tick at", t, br.Name)
		for _, timing := range br.timings {
			if timing.WeekDay == int(evaluationWeekDay)+1 &&
				timing.Hour == evaluationHour &&
				br.CurrentStatus == RelayOff &&
				timing.Minute <= evaluationMinute {
				br.setOnAndTimer(timing.OpenTimeMinutes * 60)
			}
		}
	}
}
