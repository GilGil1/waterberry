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
	baseRelay BaseRelay
}

func (gp *GPIORelay) SetConfig(id int, name string, pin uint8, timings []config.OpenTimeConfig) {
	gp.baseRelay.SetConfig(id, name, pin, timings)
}

func (gp *GPIORelay) SetFields(mode string, updateTime time.Time) error {
	gp.baseRelay.SetFields(mode, updateTime)
	return nil
}

func (gp *GPIORelay) Init() {
	gp.baseRelay.StopTime = time.Time{}

}
func (gp *GPIORelay) GetPin() uint8 {
	return gp.baseRelay.GetPin()
}

func (gp *GPIORelay) GetId() int {
	return gp.baseRelay.GetId()
}

func (gp *GPIORelay) GetName() string {
	return gp.baseRelay.GetName()
}

func (gp *GPIORelay) GetCurrentMode() string {
	return gp.baseRelay.GetCurrentMode()
}

func (gp *GPIORelay) GetSecondsOff() string {
	return gp.baseRelay.GetSecondsOff()
}

func (gp *GPIORelay) GetPropertiesMap() map[string]interface{} {
	return gp.baseRelay.GetPropertiesMap()
}

func (gp *GPIORelay) SetOn() error {
	return gp.SetMode(RelayOn)
}
func (gp *GPIORelay) SetOff() error {
	return gp.SetMode(RelayOff)
}

func (gp *GPIORelay) SetMode(mode string) error {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return err
	}
	defer rpio.Close()
	pin := rpio.Pin(gp.baseRelay.Pin)
	pin.Output()
	if mode == RelayOff {
		pin.High()
	} else {
		pin.Low()
	}
	gp.baseRelay.SetMode(mode)
	return nil
}
