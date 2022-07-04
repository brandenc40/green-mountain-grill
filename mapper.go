package gmg

import (
	"bytes"
	"encoding/binary"
)

func BytesToState(b []byte) (*State, error) {
	var m messageBody
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &m)
	if err != nil {
		return nil, err
	}
	return &State{
		WarnCode:                m.WarnCode(),
		PowerState:              m.PowerState,
		FireState:               m.FireState,
		CurveRemainTime:         m.CurveRemainTime(),
		CurrentTemperature:      m.CurrentTemperature(),
		TargetTemperature:       m.TargetTemperature(),
		Probe1Temperature:       m.Probe1Temperature(),
		Probe1TargetTemperature: m.Probe1TargetTemperature(),
		Probe2Temperature:       m.Probe2Temperature(),
		Probe2TargetTemperature: m.Probe2TargetTemperature(),
		APIVersion:              m.APIVersion,
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
	APIVersion        uint8      // 8
	FirmwareDetails   [7]uint8   // 9-15: not sure how to read this
	Probe2Temp        uint8      // 16
	Probe2TempHigh    uint8      // 17
	Probe2SetTemp     uint8      // 18
	Probe2SetTempHigh uint8      // 19
	CurveRemainTime1  uint8      // 20: not sure how to use this
	CurveRemainTime2  uint8      // 21: not sure how to use this
	CurveRemainTime3  uint8      // 22: not sure how to use this
	CurveRemainTime4  uint8      // 23: not sure how to use this
	WarnCode1         uint8      // 24
	WarnCode2         uint8      // 25
	WarnCode3         uint8      // 26
	WarnCode4         uint8      // 27
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
	return tempWithHighValue(m.GrillTemp, m.GrillTempHigh)
}

func (m messageBody) TargetTemperature() int {
	return tempWithHighValue(m.GrillSetTemp, m.GrillSetTempHigh)
}

func (m messageBody) Probe1Temperature() int {
	return tempWithHighValue(m.Probe1Temp, m.Probe1TempHigh)
}

func (m messageBody) Probe1TargetTemperature() int {
	return tempWithHighValue(m.Probe1SetTemp, m.Probe1SetTempHigh)
}

func (m messageBody) Probe2Temperature() int {
	return tempWithHighValue(m.Probe2Temp, m.Probe2TempHigh)
}

func (m messageBody) Probe2TargetTemperature() int {
	return tempWithHighValue(m.Probe2SetTemp, m.Probe2SetTempHigh)
}

func (m messageBody) CurveRemainTime() [3]int {
	return curveValue(fourByteConversion(m.CurveRemainTime1, m.CurveRemainTime2, m.CurveRemainTime3, m.CurveRemainTime4))
}

func (m messageBody) WarnCode() WarnCode {
	return WarnCode(fourByteConversion(m.WarnCode1, m.WarnCode2, m.WarnCode3, m.WarnCode4))
}
