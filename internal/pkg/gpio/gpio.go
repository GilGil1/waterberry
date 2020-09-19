package gpio

// +build !win

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
	
)

func Init() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	firstPin := rpio.Pin(17)
	firstPin.Output()
	secondPin := rpio.Pin(13)
	secondPin.Output()
	thirdPin := rpio.Pin(26)
	thirdPin.Output()
	forthPin := rpio.Pin(19)
	forthPin.Output()
	for x := 0; x < 20; x++ {
		firstPin.High()
		time.Sleep(time.Second)
		firstPin.Low()
		time.Sleep(time.Second)
		firstPin.Toggle()

		secondPin.Toggle()
		time.Sleep(time.Second)

		thirdPin.Toggle()
		time.Sleep(time.Second)

		forthPin.Toggle()
		time.Sleep(time.Second)

	}
}
