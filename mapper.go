package gmg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func BytesToState(b []byte) (*State, error) {
	var m messageBody
	fmt.Println(b)
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &m)
	if err != nil {
		return nil, err
	}
	return &State{
		WarnCode:                m.WarnCode,
		PowerState:              m.PowerState,
		FireState:               m.FireState,
		CurrentTemperature:      m.CurrentTemperature(),
		TargetTemperature:       m.TargetTemperature(),
		Probe1Temperature:       m.Probe1Temperature(),
		Probe1TargetTemperature: m.Probe1TargetTemperature(),
		Probe2Temperature:       m.Probe2Temperature(),
		Probe2TargetTemperature: m.Probe2TargetTemperature(),
	}, nil
}

// messageBody is the byte model of state response, 36 bytes are returned by the grill and
// can be mapped to this struct.
type messageBody struct {
	URPlaceholder     uint16     // 0-1: `UR`
	GrillTemp         uint8      // 2
	GrillTempHigh     uint8      // 3
	Probe1Temp        uint8      // 4
	Probe1TempHigh    uint8      // 5
	GrillSetTemp      uint8      // 6
	GrillSetTempHigh  uint8      // 7
	_                 [8]uint8   // 8-15: skip 8 bytes
	Probe2Temp        uint8      // 16
	Probe2TempHigh    uint8      // 17
	Probe2SetTemp     uint8      // 18
	Probe2SetTempHigh uint8      // 19
	CurveRemainTime   uint8      // 20: not validated
	_                 [3]uint8   // 21-23: skip 3 bytes
	WarnCode          WarnCode   // 24
	_                 [3]uint8   // 25-27: skip 3 bytes
	Probe1SetTemp     uint8      // 28
	Probe1SetTempHigh uint8      // 29
	PowerState        PowerState // 30
	GrillMode         uint8      // 31: not validated
	FireState         FireState  // 32
	FileStatePercent  uint8      // 33: not validated
	ProfileEnd        uint8      // 34: not validated
	GrillType         uint8      // 35: not validated
}

func (m messageBody) CurrentTemperature() int {
	return getTempWithHighValue(m.GrillTemp, m.GrillTempHigh)
}

func (m messageBody) TargetTemperature() int {
	return getTempWithHighValue(m.GrillSetTemp, m.GrillSetTempHigh)
}

func (m messageBody) Probe1Temperature() int {
	return getTempWithHighValue(m.Probe1Temp, m.Probe1TempHigh)
}

func (m messageBody) Probe1TargetTemperature() int {
	return getTempWithHighValue(m.Probe1SetTemp, m.Probe1SetTempHigh)
}

func (m messageBody) Probe2Temperature() int {
	return getTempWithHighValue(m.Probe2Temp, m.Probe2TempHigh)
}

func (m messageBody) Probe2TargetTemperature() int {
	return getTempWithHighValue(m.Probe2SetTemp, m.Probe2SetTempHigh)
}
