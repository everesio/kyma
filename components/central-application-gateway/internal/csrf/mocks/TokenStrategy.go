// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	apperrors "github.com/kyma-project/kyma/components/application-gateway/pkg/apperrors"

	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// TokenStrategy is an autogenerated mock type for the TokenStrategy type
type TokenStrategy struct {
	mock.Mock
}

// AddCSRFToken provides a mock function with given fields: apiRequest
func (_m *TokenStrategy) AddCSRFToken(apiRequest *http.Request) apperrors.AppError {
	ret := _m.Called(apiRequest)

	var r0 apperrors.AppError
	if rf, ok := ret.Get(0).(func(*http.Request) apperrors.AppError); ok {
		r0 = rf(apiRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(apperrors.AppError)
		}
	}

	return r0
}

// Invalidate provides a mock function with given fields:
func (_m *TokenStrategy) Invalidate() {
	_m.Called()
}
