package grillclient

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

// PowerState -
//go:generate enumer -type=PowerState -json -sql
type PowerState int

// PowerState enum values
const (
	PowerStateOff PowerState = iota
	PowerStateOn
	PowerStateFan
	PowerStateRemain
)

// FireState -
//go:generate enumer -type=FireState -json -sql
type FireState int

// FireState enum values
const (
	FireStateDefault FireState = iota
	FireStateOff
	FireStateStartup
	FireStateRunning
	FireStateCoolDown
	FireStateFail
)

// WarnCode -
//go:generate enumer -type=WarnCode -json -sql
type WarnCode int

// WarnCode enum values
const (
	WarnCodeNone      WarnCode = iota
	WarnCodeLowPellet WarnCode = 128
)
