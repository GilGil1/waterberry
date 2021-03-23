// +build linux

package gpio

import (
	"fmt"
	"time"
	config "waterberry/internal/pkg/config"

	"github.com/stianeikeland/go-rpio"
)

const (
	PIN_HIGH = "high"
	PIN_LOW  = "low"
)

type GPIORelay struct {
	base *BaseRelay
}

func NewGPIORelay() *GPIORelay {
	relay := GPIORelay{}
	relay.base = &BaseRelay{}
	return &relay
}

func (gp *GPIORelay) Init() error {
	gp.SetOff()
	if len(gp.base.timings) > 0 {
		go StartTimers(gp, gp.base)
	}
	return nil
}

func (gp *GPIORelay) SetConfig(id int, name string, pin uint8, timings []config.OpenTimeConfig) {
	gp.base.SetConfig(id, name, pin, timings)
}

func (gp *GPIORelay) GetPropertiesMap() map[string]interface{} {
	return gp.base.GetPropertiesMap()
}

func (gp *GPIORelay) SetOn(stopMinutes int) error {
	go SetTimeOff(gp, stopMinutes)
	gp.base.StopTime = time.Now().Add(time.Duration(stopMinutes) * time.Minute)
	return gp.setMode(RelayOn)
}
func (gp *GPIORelay) SetOff() error {
	return gp.setMode(RelayOff)
}

func (gp *GPIORelay) setMode(mode string) error {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return err
	}
	defer rpio.Close()
	pin := rpio.Pin(gp.base.Pin)
	pin.Output()
	if mode == RelayOff {
		pin.High()
		gp.base.SetOff()
	} else {
		pin.Low()
		gp.base.SetOn()
	}
	return nil
}
