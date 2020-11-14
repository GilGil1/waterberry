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
	baseRelay BaseRelay
}

func (gp *StubRelay) SetConfig(id int, name string, pin uint8, timings []config.OpenTimeConfig) {
	gp.baseRelay.SetConfig(id, name, pin, timings)
}

func (gp *StubRelay) SetFields(mode string, updateTime time.Time) error {
	gp.baseRelay.SetFields(mode, updateTime)
	return nil
}

func (gp *StubRelay) Init() {

}

func (gp *StubRelay) GetPin() uint8 {
	return gp.baseRelay.Pin
}

func (gp *StubRelay) GetName() string {
	return gp.baseRelay.Name
}

func (gp *StubRelay) GetCurrentMode() string {
	return gp.baseRelay.CurrentStatus
}

func (gp StubRelay) GetId() int {
	return gp.baseRelay.ID
}

func (gp StubRelay) GetSecondsOff() string {
	return "A lot"
}

func (gp *StubRelay) SetOn() error {
	return gp.SetMode(RelayOn)
}
func (gp *StubRelay) SetOff() error {
	return gp.SetMode(RelayOff)
}

func (gp *StubRelay) SetMode(mode string) error {
	gp.baseRelay.SetMode(mode)

	return nil
}

func (gp *StubRelay) GetPropertiesMap() map[string]interface{} {
	return gp.baseRelay.GetPropertiesMap()

}
