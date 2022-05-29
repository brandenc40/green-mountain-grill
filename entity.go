package gmg

// State - The current state of the grill
type State struct {
	CurrentTemperature      int        `json:"current_temperature"`
	TargetTemperature       int        `json:"target_temperature"`
	Probe1Temperature       int        `json:"probe1_temperature"`
	Probe1TargetTemperature int        `json:"probe1_target_temperature"`
	Probe2Temperature       int        `json:"probe2_temperature"`
	Probe2TargetTemperature int        `json:"probe2_target_temperature"`
	WarnCode                WarnCode   `json:"warn_code"`
	PowerState              PowerState `json:"power_state"`
	FireState               FireState  `json:"fire_state"`
}

// IsOn - true if the grill is turned on
func (s *State) IsOn() bool {
	return s.PowerState == PowerStateOn || s.PowerState == PowerStateColdSmoke
}

// PowerState -
//go:generate go run github.com/alvaroloes/enumer -type=PowerState -json -sql
type PowerState uint8

// PowerState enum values
const (
	PowerStateOff PowerState = iota
	PowerStateOn
	PowerStateFan
	PowerStateColdSmoke
)

// FireState -
//go:generate go run github.com/alvaroloes/enumer -type=FireState -json -sql
type FireState uint8

// FireState enum values
const (
	FireStateDefault   FireState = 0
	FireStateOff       FireState = 1
	FireStateStartup   FireState = 2
	FireStateRunning   FireState = 3
	FireStateCoolDown  FireState = 4
	FireStateFail      FireState = 5
	FireStateColdSmoke FireState = 198
)

// WarnCode -
//go:generate go run github.com/alvaroloes/enumer -type=WarnCode -json -sql
type WarnCode uint8

// WarnCode enum values
// TODO: VALIDATE, CURRENTLY NOT VALIDATED FOR ALL ERROR CODES
const (
	WarnCodeNone               WarnCode = 0
	WarnCodeFanMotorOverload   WarnCode = 1
	WarnCodeAugerMotorOverload WarnCode = 2
	WarnCodeLowVoltage         WarnCode = 4
	WarnCodeIgniterOverload    WarnCode = 8
	WarnCodeLowPellet          WarnCode = 128
)
