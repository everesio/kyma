// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	applications "github.com/kyma-project/kyma/components/application-gateway/internal/metadata/applications"
	apperrors "github.com/kyma-project/kyma/components/application-gateway/pkg/apperrors"

	mock "github.com/stretchr/testify/mock"
)

// ServiceRepository is an autogenerated mock type for the ServiceRepository type
type ServiceRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: appName, id
func (_m *ServiceRepository) Get(appName string, id string) (applications.Service, apperrors.AppError) {
	ret := _m.Called(appName, id)

	var r0 applications.Service
	if rf, ok := ret.Get(0).(func(string, string) applications.Service); ok {
		r0 = rf(appName, id)
	} else {
		r0 = ret.Get(0).(applications.Service)
	}

	var r1 apperrors.AppError
	if rf, ok := ret.Get(1).(func(string, string) apperrors.AppError); ok {
		r1 = rf(appName, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(apperrors.AppError)
		}
	}

	return r0, r1
}
