# Green Mountain Grill

### Observe and Control your grill with Go

The `grillclient` package can be used as a universal client for 
interacting with your Green Mountain Grill. 

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

// Params - Parameters to build a new Client
type Params struct {
	GrillIP         net.IP
	GrillPort       int
	Logger          *logrus.Logger
	ReadTimeout     time.Duration // default 2 seconds
	WriteTimeout    time.Duration // default 1 second
	MaxConnAttempts int           // default 5
}
```

Building a new client
```go
import "github.com/brandenc40/gmg/grillclient"

params := grillclient.Params{
	GrillIP:   "192.168.1.2", // example, this will change
	GrillPort: "8080", // this should be the same for all grills... I think...
}
client := grillclient.New(params)
```

Planned features to add:
- track temp over time
- alerts for when temps are reached

__Note: this was tested on my grill which is a Daniel Boone Prime purchased 
in 2021. I'm  not sure if this will work properly on other models.__

### Still a work in progress so feel free to assist with building this codebase, any help would be appreciated

## Grill State Data Parse

> Shout out to https://github.com/Aenima4six2/gmg and https://github.com/FeatherKing/grillsrv 
> for doing a lot of the leg work on figuring out the commands to send and the 
> data returned by the grill.

```
EXAMPLE: GRILL OFF
INDEX:  0  1  2  3 4  5 6   7 8 9  10 11 12 13 14 15 16 17 18 19 20  21  22  23  24 25 26 27 28 29 30 31 32 33 34 35
BYTES: [85 82 97 0 89 2 150 0 5 11 20 50 25 25 25 25 89 2  0  0  255 255 255 255 0  0  0  0  0  0  0  0  1  0  0  3 ]

VALUE INDICIES
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
```
