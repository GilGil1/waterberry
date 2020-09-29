// +build linux

package gpio

import (
	"fmt"

	"github.com/stianeikeland/go-rpio"
)

type GPIORelay struct {
	Pin           uint8
	Name          string
	ID            int
	CurrentStatus string
}

func (gp *GPIORelay) SetConfig(id int, name string, pin uint8) {
	gp.ID = id
	gp.Name = name
	gp.Pin = pin
}
func (gp *GPIORelay) Init() {

}
func (gp *GPIORelay) GetPin() uint8 {
	return gp.Pin
}

func (gp *GPIORelay) GetId() int {
	return gp.ID
}

func (gp *GPIORelay) GetName() string {
	return gp.Name
}

func (gp *GPIORelay) SetOn() error {
	fmt.Printf("Before On relay # %d (pin = %d)", gp.ID, gp.Pin)

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return err
	}
	defer rpio.Close()
	pin := rpio.Pin(gp.Pin)
	pin.Output()
	pin.High()
	gp.CurrentStatus = RelayOn
	fmt.Printf("After On relay # %d (pin = %d)", gp.ID, gp.Pin)

	return nil
}
func (gp *GPIORelay) SetOff() error {
	fmt.Printf("Before off of relay # %d (pin = %d)", gp.ID, gp.Pin)
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return err
	}
	defer rpio.Close()
	pin := rpio.Pin(gp.Pin)
	pin.Output()
	pin.Low()
	gp.CurrentStatus = RelayOff
	fmt.Printf("After off of relay # %d (pin = %d)", gp.ID, gp.Pin)
	return nil
}

func (gp *GPIORelay) GetCurrentMode() string {
	return gp.CurrentStatus
}
