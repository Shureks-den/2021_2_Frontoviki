// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import models "yula/internal/models"

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// Insert provides a mock function with given fields: _a0
func (_m *UserRepository) Insert(_a0 *models.UserData) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.UserData) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SelectByEmail provides a mock function with given fields: email
func (_m *UserRepository) SelectByEmail(email string) (*models.UserData, error) {
	ret := _m.Called(email)

	var r0 *models.UserData
	if rf, ok := ret.Get(0).(func(string) *models.UserData); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SelectById provides a mock function with given fields: userId
func (_m *UserRepository) SelectById(userId int64) (*models.UserData, error) {
	ret := _m.Called(userId)

	var r0 *models.UserData
	if rf, ok := ret.Get(0).(func(int64) *models.UserData); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *UserRepository) Update(_a0 *models.UserData) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.UserData) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}