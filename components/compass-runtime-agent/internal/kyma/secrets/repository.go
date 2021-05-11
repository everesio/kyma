// Package secrets contains components for accessing/modifying client secrets
package secrets

import (
	"context"

	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/apperrors"
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/k8sconsts"
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/kyma/secrets/strategy"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Repository contains operations for managing client credentials
//go:generate mockery --name Repository
type Repository interface {
	Create(application string, appUID types.UID, name, serviceID string, data strategy.SecretData) apperrors.AppError
	Get(name string) (strategy.SecretData, apperrors.AppError)
	Delete(name string) apperrors.AppError
	Upsert(application string, appUID types.UID, name, secretID string, data strategy.SecretData) apperrors.AppError
}

type repository struct {
	secretsManager Manager
}

// Manager contains operations for managing k8s secrets
//go:generate mockery --name Manager
type Manager interface {
	Create(ctx context.Context, secret *v1.Secret, options metav1.CreateOptions) (*v1.Secret, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1.Secret, error)
	Delete(ctx context.Context, name string, options metav1.DeleteOptions) error
	Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) (*v1.Secret, error)
}

// NewRepository creates a new secrets repository
func NewRepository(secretsManager Manager) Repository {
	return &repository{
		secretsManager: secretsManager,
	}
}

// Create adds a new secret with one entry containing specified clientId and clientSecret
func (r *repository) Create(application string, appUID types.UID, name, serviceID string, data strategy.SecretData) apperrors.AppError {
	secret := makeSecret(name, serviceID, application, appUID, data)
	return r.create(application, secret, name)
}

func (r *repository) Get(name string) (data strategy.SecretData, error apperrors.AppError) {
	secret, err := r.secretsManager.Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return strategy.SecretData{}, apperrors.NotFound("Secret %s not found", name)
		}
		return strategy.SecretData{}, apperrors.Internal("Getting %s secret failed, %s", name, err.Error())
	}

	return secret.Data, nil
}

func (r *repository) Delete(name string) apperrors.AppError {
	err := r.secretsManager.Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return apperrors.Internal("Deleting %s secret failed, %s", name, err.Error())
	}
	return nil
}

func (r *repository) Upsert(application string, appUID types.UID, name, serviceID string, data strategy.SecretData) apperrors.AppError {
	secret := makeSecret(name, serviceID, application, appUID, data)

	_, err := r.secretsManager.Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return r.create(application, secret, name)
		}
		return apperrors.Internal("Updating %s secret failed, %s", name, err.Error())
	}
	return nil
}

func (r *repository) create(application string, secret *v1.Secret, name string) apperrors.AppError {
	_, err := r.secretsManager.Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return apperrors.AlreadyExists("Secret %s already exists", name)
		}
		return apperrors.Internal("Creating %s secret failed, %s", name, err.Error())
	}
	return nil
}

func makeSecret(name, serviceID, application string, appUID types.UID, data strategy.SecretData) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				k8sconsts.LabelApplication: application,
				k8sconsts.LabelServiceId:   serviceID,
			},
			OwnerReferences: k8sconsts.CreateOwnerReferenceForApplication(application, appUID),
		},
		Data: data,
	}
}
