# Green Mountain Grill

### Observe and Control your grill with Go

__Note: this was tested on my grill which is a Daniel Boone Prime purchased 
in 2021. I'm  not sure if this will work properly on other models.__

### Frontend stuff is still a work in progress so feel free to assist with building this codebase.

## Grill Client

```go
// Client - Green Mountain Grill client interface definition
type Client interface {
	IsAvailable() bool
	GetState() (*State, error)
	GetID() ([]byte, error)
	GetFirmware() ([]byte, error)
	SetGrillTemp(temp int) error
	SetProbe1Target(temp int) error
	SetProbe2Target(temp int) error
	PowerOn() error
	PowerOnColdSmoke() error
	PowerOff() error
}
```

Building a new client
```go
import gmg "github.com/brandenc40/green-mountain-grill"

client := gmg.New(
	net.ParseIP("192.168.1.2"), // example, this will change
	8080, // this should be the same for all grills... I think...
	gmg.WithZapLogger(p.Logger),
)
```

## Web Server

WORK IN PROGRESS

Planned features to add:
- track temp over time
- alerts for when temps are reached

## Grill State Data Parse

> Shout out to https://github.com/Aenima4six2/gmg and https://github.com/FeatherKing/grillsrv 
> for doing a lot of the leg work on figuring out the commands to send and the 
> data returned by the grill.

```go
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
```
