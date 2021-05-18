package metadata

import (
	"testing"
)

func TestServiceDefinitionService_GetAPI(t *testing.T) {

	//t.Run("should get API", func(t *testing.T) {
	//	// given
	//	applicationServiceAPI := &applications.ServiceAPI{}
	//	applicationService := applications.Service{API: applicationServiceAPI}
	//	serviceAPI := &model.API{}
	//
	//	serviceRepository := new(applicationmocks.ServiceRepository)
	//	serviceRepository.On("Get", "app1", "service-1", "api-1").Return(applicationService, nil)
	//
	//	serviceAPIService := new(serviceapimocks.Service)
	//	serviceAPIService.On("Read", applicationServiceAPI).Return(serviceAPI, nil)
	//
	//	service := NewServiceDefinitionService(serviceAPIService, serviceRepository)
	//
	//	// when
	//	result, err := service.GetAPI("app1", "service-1", "api-1")
	//
	//	// then
	//	require.NoError(t, err)
	//
	//	assert.Equal(t, serviceAPI, result)
	//})
	//
	//t.Run("should return not found error if service does not exist", func(t *testing.T) {
	//	// given
	//	serviceRepository := new(applicationmocks.ServiceRepository)
	//	serviceRepository.On("Get", "app1", "service-1", "api-1").Return(applications.Service{}, apperrors.NotFound("missing"))
	//
	//	service := NewServiceDefinitionService(nil, serviceRepository)
	//
	//	// when
	//	result, err := service.GetAPI("app1", "service-1", "api-1")
	//
	//	// then
	//	assert.Error(t, err)
	//	assert.Nil(t, result)
	//	assert.Equal(t, apperrors.CodeNotFound, err.Code())
	//})
	//
	//t.Run("should return internal error if service does not exist", func(t *testing.T) {
	//	// given
	//	serviceRepository := new(applicationmocks.ServiceRepository)
	//	serviceRepository.On("Get", "app1", "service-1", "api-1").Return(applications.Service{}, apperrors.Internal("some error"))
	//
	//	service := NewServiceDefinitionService(nil, serviceRepository)
	//
	//	// when
	//	result, err := service.GetAPI("app1", "service-1", "api-1")
	//
	//	// then
	//	assert.Error(t, err)
	//	assert.Nil(t, result)
	//	assert.Equal(t, apperrors.CodeInternal, err.Code())
	//	assert.Contains(t, err.Error(), "some error")
	//})
	//
	//t.Run("should return bad request if service does not have API", func(t *testing.T) {
	//	// given
	//	serviceRepository := new(applicationmocks.ServiceRepository)
	//	serviceRepository.On("Get", "app1", "service-1", "api-1").Return(applications.Service{}, nil)
	//
	//	service := NewServiceDefinitionService(nil, serviceRepository)
	//
	//	// when
	//	result, err := service.GetAPI("app1", "service-1", "api-1")
	//
	//	// then
	//	assert.Error(t, err)
	//	assert.Nil(t, result)
	//	assert.Equal(t, apperrors.CodeWrongInput, err.Code())
	//})
	//
	//t.Run("should return internal error if reading service API fails", func(t *testing.T) {
	//	// given
	//	applicationServiceAPI := &applications.ServiceAPI{}
	//	applicationService := applications.Service{API: applicationServiceAPI}
	//
	//	serviceRepository := new(applicationmocks.ServiceRepository)
	//	serviceRepository.On("Get", "app1", "service-1", "api-1").Return(applicationService, nil)
	//
	//	serviceAPIService := new(serviceapimocks.Service)
	//	serviceAPIService.On("Read", applicationServiceAPI).Return(nil, apperrors.Internal("some error"))
	//
	//	service := NewServiceDefinitionService(serviceAPIService, serviceRepository)
	//
	//	// when
	//	result, err := service.GetAPI("app1", "service-1", "api-1")
	//
	//	// then
	//	assert.Error(t, err)
	//	assert.Nil(t, result)
	//	assert.Equal(t, apperrors.CodeInternal, err.Code())
	//	assert.Contains(t, err.Error(), "some error")
	//})
}