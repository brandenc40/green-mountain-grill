// Code generated by mockery v2.8.0. DO NOT EDIT.

package mocks

import (
	model "github.com/brandenc40/green-mountain-grill/server/respository/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// GetStateHistory provides a mock function with given fields: sessionUUID
func (_m *Repository) GetStateHistory(sessionUUID uuid.UUID) ([]*model.GrillState, error) {
	ret := _m.Called(sessionUUID)

	var r0 []*model.GrillState
	if rf, ok := ret.Get(0).(func(uuid.UUID) []*model.GrillState); ok {
		r0 = rf(sessionUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.GrillState)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(sessionUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertStateData provides a mock function with given fields: state
func (_m *Repository) InsertStateData(state *model.GrillState) error {
	ret := _m.Called(state)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.GrillState) error); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
