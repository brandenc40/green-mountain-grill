package mapper

import (
	"github.com/brandenc40/green-mountain-grill/client"
	"github.com/brandenc40/green-mountain-grill/internal/respository/model"
	"github.com/google/uuid"
)

// GrillStateEntityToModel -
func GrillStateEntityToModel(gs *client.State, sessionUUID uuid.UUID) *model.GrillState {
	return &model.GrillState{
		SessionUUID:             sessionUUID,
		CurrentTemperature:      gs.CurrentTemperature,
		TargetTemperature:       gs.TargetTemperature,
		Probe1Temperature:       gs.Probe1Temperature,
		Probe1TargetTemperature: gs.Probe1TargetTemperature,
		Probe2Temperature:       gs.Probe2Temperature,
		Probe2TargetTemperature: gs.Probe2TargetTemperature,
		WarnCode:                gs.WarnCode,
		PowerState:              gs.PowerState,
		FireState:               gs.FireState,
	}
}
