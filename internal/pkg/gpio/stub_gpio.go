package gpio

// Enumerations to define the relay status
const (
	RelayOff    = "off"
	RelayOn     = "on"
	RelayToggle = "toggle"
)

// IRelay ...
// An interface to a genericv relay.
// It may be a stub or a real io
type IRelay interface {
	SetConfig(id int, name string, pin uint8)
	Init()
	GetPin() uint8
	GetId() int
	GetName() string
	GetSecondsOff() string
	GetCurrentMode() string
	SetOn() error
	SetOff() error
}

// var globalConfig config.GlobalConfig
// var RelaysInfo [4]RelayInfo

// type RelayInfo struct {
// 	CurrentMode string
// 	Config      config.RelayConfig
// }

// StubRelay ...
// This is a stub of GPIORelet
type StubRelay struct {
	ID            int
	Pin           uint8
	Name          string
	CurrentStatus string
}

func (gp *StubRelay) SetConfig(id int, name string, pin uint8) {
	gp.ID = id
	gp.Name = name
	gp.Pin = pin
}
func (gp *StubRelay) Init() {

}
func (gp *StubRelay) GetPin() uint8 {
	return gp.Pin
}

func (gp *StubRelay) GetName() string {
	return gp.Name
}

func (gp *StubRelay) GetCurrentMode() string {
	return gp.CurrentStatus
}

func (gp *StubRelay) GetId() int {
	return gp.ID
}

func (gp *StubRelay) GetSecondsOff() string {
	return "A lot"
}

func (gp *StubRelay) SetOn() error {
	gp.CurrentStatus = RelayOn
	return nil
}
func (gp *StubRelay) SetOff() error {
	gp.CurrentStatus = RelayOff
	return nil
}

func GetPropertiesMap(relay IRelay) map[string]interface{} {
	var singleMap = make(map[string]interface{})
	singleMap["id"] = relay.GetId()
	singleMap["name"] = relay.GetName()
	singleMap["pin"] = relay.GetPin()
	singleMap["mode"] = relay.GetCurrentMode()
	singleMap["seconds_off"] = relay.GetSecondsOff()
	return singleMap

}
