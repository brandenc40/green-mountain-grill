package repository

import (
	"errors"

	"github.com/google/uuid"

	"github.com/brandenc40/gmg/internal/respository/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var models = []interface{}{
	&model.GrillState{},
}

type Repository interface {
	GetStateHistory(sessionUUID uuid.UUID) ([]*model.GrillState, error)
	InsertStateData(state *model.GrillState) error
}

func New() (Repository, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"))
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(models...); err != nil {
		return nil, err
	}
	return &repository{db: db}, nil
}

type repository struct {
	db *gorm.DB
}

func (r *repository) GetStateHistory(sessionUUID uuid.UUID) ([]*model.GrillState, error) {
	var stateHist []*model.GrillState
	err := r.db.Where("session_uuid = ?", sessionUUID).Find(&stateHist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*model.GrillState{}, nil
		}
		return nil, err
	}
	return stateHist, nil
}

func (r *repository) InsertStateData(state *model.GrillState) error {
	return r.db.Create(&state).Error
}