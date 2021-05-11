package secrets

import (
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/apperrors"
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/k8sconsts"
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/kyma/applications"
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/kyma/model"
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/kyma/secrets/strategy"
	"k8s.io/apimachinery/pkg/types"
)

type modificationFunction func(modStrategy strategy.ModificationStrategy, application string, appUID types.UID, name, serviceID string, newData strategy.SecretData) apperrors.AppError

//go:generate mockery --name Service
type Service interface {
	Get(application string, credentials applications.Credentials) (model.Credentials, apperrors.AppError)
	Create(application string, appUID types.UID, serviceID string, credentials *model.Credentials) (applications.Credentials, apperrors.AppError)
	Upsert(application string, appUID types.UID, serviceID string, credentials *model.Credentials) (applications.Credentials, apperrors.AppError)
	Delete(name string) apperrors.AppError
}

type service struct {
	repository      Repository
	strategyFactory strategy.Factory
}

func NewService(repository Repository, strategyFactory strategy.Factory) Service {
	return &service{
		repository:      repository,
		strategyFactory: strategyFactory,
	}
}

func (s *service) Create(application string, appUID types.UID, serviceID string, credentials *model.Credentials) (applications.Credentials, apperrors.AppError) {
	return s.modifySecret(application, appUID, serviceID, credentials, s.createSecret)
}

func (s *service) Get(application string, credentials applications.Credentials) (model.Credentials, apperrors.AppError) {
	accessStrategy, err := s.strategyFactory.NewSecretAccessStrategy(&credentials)
	if err != nil {
		return model.Credentials{}, err.Append("Failed to initialize strategy")
	}

	data, err := s.repository.Get(credentials.SecretName)
	if err != nil {
		return model.Credentials{}, err
	}

	return accessStrategy.ToCredentials(data, &credentials)
}

func (s *service) Upsert(application string, appUID types.UID, serviceID string, credentials *model.Credentials) (applications.Credentials, apperrors.AppError) {
	return s.modifySecret(application, appUID, serviceID, credentials, s.upsertSecret)
}

func (s *service) Delete(name string) apperrors.AppError {
	return s.repository.Delete(name)
}

func (s *service) modifySecret(application string, appUID types.UID, serviceID string, credentials *model.Credentials, modFunction modificationFunction) (applications.Credentials, apperrors.AppError) {
	if credentials == nil {
		return applications.Credentials{}, nil
	}

	modStrategy, err := s.strategyFactory.NewSecretModificationStrategy(credentials)
	if err != nil {
		return applications.Credentials{}, err.Append("Failed to initialize strategy")
	}

	if !modStrategy.CredentialsProvided(credentials) {
		return applications.Credentials{}, nil
	}

	name := k8sconsts.GetResourceName(application, serviceID)

	secretData, err := modStrategy.CreateSecretData(credentials)
	if err != nil {
		return applications.Credentials{}, err.Append("Failed to create secret data")
	}

	err = modFunction(modStrategy, application, appUID, name, serviceID, secretData)
	if err != nil {
		return applications.Credentials{}, err
	}

	return modStrategy.ToCredentialsInfo(credentials, name), nil
}

func (s *service) upsertSecret(modStrategy strategy.ModificationStrategy, application string, appUID types.UID, name, serviceID string, newData strategy.SecretData) apperrors.AppError {
	currentData, err := s.repository.Get(name)
	if err != nil {
		if err.Code() == apperrors.CodeNotFound {
			return s.repository.Create(application, appUID, name, serviceID, newData)
		}

		return err
	}

	if modStrategy.ShouldUpdate(currentData, newData) {
		return s.repository.Upsert(application, appUID, name, serviceID, newData)
	}

	return nil
}

func (s *service) createSecret(_ strategy.ModificationStrategy, application string, appUID types.UID, name, serviceID string, newData strategy.SecretData) apperrors.AppError {
	return s.repository.Create(application, appUID, name, serviceID, newData)
}
