package grillclient

import "fmt"

// Command -
type Command string

// Commands that are accepted by the Green Mountain Grill
const (
	CommandGetInfo               Command = "UR001!"
	CommandSetGrillTemp          Command = "UT%d!"
	CommandSetProbeTemp          Command = "UF%d!"
	CommandPowerOn               Command = "UK001!"
	CommandPowerOff              Command = "UK004!"
	CommandGetGrillID            Command = "UL!"
	CommandGetGrillFirmware      Command = "UN!"
	CommandBroadcastToClientMode Command = "UH%c%c%s%c%s!"
)

// Bytes -
func (c Command) Bytes(args ...interface{}) []byte {
	if len(args) > 0 {
		return []byte(fmt.Sprintf(string(c), args...))
	}
	return []byte(c)
}
