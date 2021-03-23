package gpio

import (
	"time"
	config "waterberry/internal/pkg/config"
)

// Enumerations to define the relay status
const (
	RelayOff    = "off"
	RelayOn     = "on"
	RelayToggle = "toggle"
)

// IRelay ...
// An interface to a genericv relay.
// It may be a stub or a real io

type StubRelay struct {
	base *BaseRelay
}

func NewStubRelay() *StubRelay {
	relay := StubRelay{}
	relay.base = &BaseRelay{}
	return &relay
}

func (gp *StubRelay) SetConfig(id int, name string, pin uint8, timings []config.OpenTimeConfig) {
	gp.base.SetConfig(id, name, pin, timings)
}

func (gp *StubRelay) GetBaseRelay() *BaseRelay {
	return gp.base
}

func (gp *StubRelay) SetOn(stopMinutes int) error {
	go SetTimeOff(gp, stopMinutes)
	gp.base.StopTime = time.Now().Add(time.Duration(stopMinutes) * time.Minute)

	return gp.base.SetOn()
}
func (gp *StubRelay) SetOff() error {
	return gp.base.SetOff()
}

func (gp *StubRelay) Init() error {
	gp.base.SetOff()
	if len(gp.base.timings) > 0 {
		go StartTimers(gp, gp.base)
	}
	return nil
}

// func (gp *StubRelay) SetMode(mode string) error {
// 	return gp.base.SetMode(mode)
// }

func (gp *StubRelay) GetPropertiesMap() map[string]interface{} {
	return gp.base.GetPropertiesMap()

}
