package client

// State Bytes Locations
// UR[2 Byte Grill Temp][2 Byte food probe Temp][2 Byte Target Temp][skip 22 bytes][2 Byte target food probe][1byte on/off/fan][5 byte tail]
const (
	grillTemp         = 2
	grillTempHigh     = 3
	probeTemp         = 4
	probeTempHigh     = 5
	grillSetTemp      = 6
	grillSetTempHigh  = 7
	probe2Temp        = 16
	probe2TempHigh    = 17
	probe2SetTemp     = 18
	probe2SetTempHigh = 19
	curveRemainTime   = 20 // not validated
	warnCode          = 24
	probeSetTemp      = 28
	probeSetTempHigh  = 29
	powerState        = 30
	grillMode         = 31 // not validated
	fireState         = 32
	fileStatePercent  = 33 // not validated
	profileEnd        = 34 // not validated
	grillType         = 35 // not validated
)

// GetStateResponseToState -
func GetStateResponseToState(response []byte) *State {
	state := &State{
		WarnCode:                WarnCode(response[warnCode]),
		PowerState:              PowerState(response[powerState]),
		FireState:               FireState(response[fireState]),
		CurrentTemperature:      getTempWithHighVal(response, grillTemp, grillTempHigh),
		TargetTemperature:       getTempWithHighVal(response, grillSetTemp, grillSetTempHigh),
		Probe1Temperature:       getTempWithHighVal(response, probeTemp, probeTempHigh),
		Probe1TargetTemperature: getTempWithHighVal(response, probeSetTemp, probeSetTempHigh),
		Probe2Temperature:       getTempWithHighVal(response, probe2Temp, probe2TempHigh),
		Probe2TargetTemperature: getTempWithHighVal(response, probe2SetTemp, probe2SetTempHigh),
	}
	return state
}

func getTempWithHighVal(data []byte, tmpIdx, highIdx int) int {
	high := int(data[highIdx])
	// a high value of 2 represents that the temp is not available
	if high == 2 {
		return 0
	}
	return int(data[tmpIdx]) + high*256
}
