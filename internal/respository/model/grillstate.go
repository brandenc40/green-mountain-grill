package model

import (
	"time"

	"github.com/brandenc40/green-mountain-grill/grillclient"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GrillState - The current state of the grill
type GrillState struct {
	ID                      uint      `gorm:"primarykey"`
	CreatedAt               time.Time `gorm:"index"`
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt `gorm:"index"`
	SessionUUID             uuid.UUID      `gorm:"index"`
	CurrentTemperature      int
	TargetTemperature       int
	Probe1Temperature       int
	Probe1TargetTemperature int
	Probe2Temperature       int
	Probe2TargetTemperature int
	WarnCode                grillclient.WarnCode
	PowerState              grillclient.PowerState
	FireState               grillclient.FireState
}
