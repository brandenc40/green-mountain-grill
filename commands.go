package gmg

import "fmt"

// Command -
type Command string

// Commands that are accepted by the Green Mountain Grill
const (
	CommandGetInfo          Command = "UR001!"
	CommandGetGrillID       Command = "UL!"
	CommandGetGrillFirmware Command = "UN!"
	CommandSetGrillTemp     Command = "UT%03d!"
	CommandSetProbe1Temp    Command = "UF%03d!"
	CommandSetProbe2Temp    Command = "Uf%03d!"
	CommandPowerOn          Command = "UK001!"
	CommandPowerOnColdSmoke Command = "UK002!"
	CommandPowerOff         Command = "UK004!"
)

// Build -
func (c Command) Build(args ...interface{}) []byte {
	if len(args) > 0 {
		formatted := fmt.Sprintf(string(c), args...)
		return []byte(formatted)
	}
	return []byte(c)
}
