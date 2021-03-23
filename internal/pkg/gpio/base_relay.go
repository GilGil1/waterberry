package gpio

import (
	"fmt"
	"time"
	config "waterberry/internal/pkg/config"
)

type IRelay interface {
	SetConfig(id int, name string, pin uint8, timings []config.OpenTimeConfig)
	SetOn(OpenMinutes int) error
	SetOff() error
	GetPropertiesMap() map[string]interface{}
	Init() error
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
		return fmt.Errorf("mode %s unsupported", mode)
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

func (br *BaseRelay) SetConfig(id int, name string, pin uint8, timings []config.OpenTimeConfig) {
	br.ID = id
	br.Name = name
	br.Pin = pin
	br.timings = timings
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
	if br.StopTime.Before(time.Now()) {
		if len(br.timings) > 0 {
			return "Timers Exist"
		} else {
			return "No Timers Exist on Relay"

		}
	} else {
		return fmt.Sprintf("%4.0f", time.Until(br.StopTime).Seconds())

	}
}

func (br *BaseRelay) GetCurrentMode() string {
	return br.CurrentStatus
}

func (br *BaseRelay) SetOn() error {

	return br.setMode(RelayOn)
}

func (br *BaseRelay) SetOff() error {
	return br.setMode(RelayOff)

}

func (br *BaseRelay) setMode(mode string) error {
	if mode == RelayOn {
		br.CurrentStatus = mode

		printModeChange(*br)
	} else if mode == RelayOff {
		br.CurrentStatus = mode
		br.StopTime = time.Now()
		printModeChange(*br)
	} else {
		return fmt.Errorf("mode %s unsupported", mode)
	}
	return nil
}
func SetTimeOff(relay IRelay, minutes int) {
	fmt.Printf("Set Timer in %d minutes \n", minutes)
	time.Sleep(time.Duration(int64(minutes) * int64(time.Minute)))
	relay.SetOff()
	fmt.Printf("Relay Set Off \n")
}

func StartTimers(relay IRelay, base *BaseRelay) {
	ticker := time.NewTicker(time.Millisecond * 10000)
	// Evaluation Loop
	for t := range ticker.C {
		evalutationTime := time.Now()
		evaluationWeekDay := evalutationTime.Weekday()
		evaluationHour := evalutationTime.Hour()
		evaluationMinute := evalutationTime.Minute()
		fmt.Printf("Evaluation Cycle at:%v Relay:%s, status:%s, Stop time:%v \n",
			t,
			base.Name,
			base.CurrentStatus,
			base.StopTime.Format(time.RFC3339))
		for _, timing := range base.timings {
			if timing.WeekDay == int(evaluationWeekDay)+1 && timing.Hour == evaluationHour && base.CurrentStatus == RelayOff && timing.Minute == evaluationMinute {
				relay.SetOn(timing.OpenTimeMinutes)

			}
		}
	}
}
