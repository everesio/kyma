package applications

import (
	"fmt"
	"github.com/kyma-project/kyma/components/application-operator/pkg/apis/applicationconnector/v1alpha1"
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/kyma/model"
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/kyma/secrets"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultDescription = "Description not provided"

const (
	connectedApp = "connected-app"
)

//go:generate mockery --name=Converter
type Converter interface {
	Do(application model.Application) v1alpha1.Application
}

type converter struct {
}

func NewConverter() Converter {
	return converter{}
}

func (c converter) Do(application model.Application) v1alpha1.Application {

	prepareLabels := func(directorLabels model.Labels) map[string]string {
		labels := make(map[string]string)

		labels[connectedApp] = application.Name

		return labels
	}

	description := application.Description
	if description == "" {
		description = defaultDescription
	}

	return v1alpha1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "applicationconnector.kyma-project.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: application.Name,
		},
		Spec: v1alpha1.ApplicationSpec{
			Description:      description,
			SkipInstallation: false,
			Labels:           prepareLabels(application.Labels),
			Services:         c.toServices(application.Name, application.APIPackages),
			CompassMetadata:  c.toCompassMetadata(application.ID, application.SystemAuthsIDs),
		},
	}
}

func (c converter) toServices(applicationName string, packages []model.APIPackage) []v1alpha1.Service {
	services := make([]v1alpha1.Service, 0, len(packages))

	for _, p := range packages {
		services = append(services, c.toService(applicationName, p))
	}

	return services
}

func (c converter) toService(applicationName string, apiPackage model.APIPackage) v1alpha1.Service {

	description := defaultDescription
	if apiPackage.Description != nil && *apiPackage.Description != "" {
		description = *apiPackage.Description
	}

	return v1alpha1.Service{
		ID:                        apiPackage.ID,
		Identifier:                "", // not available in the Director's API
		Name:                      createServiceName(apiPackage.Name, apiPackage.ID),
		AuthCreateParameterSchema: apiPackage.InstanceAuthRequestInputSchema,
		DisplayName:               apiPackage.Name,
		Description:               description,
		Entries:                   c.toServiceEntries(applicationName, apiPackage),
	}
}

func (c converter) toServiceEntries(applicationName string, apiPackage model.APIPackage) []v1alpha1.Entry {
	entries := make([]v1alpha1.Entry, 0, len(apiPackage.APIDefinitions)+len(apiPackage.EventDefinitions))

	for _, apiDefinition := range apiPackage.APIDefinitions {
		entries = append(entries, c.toAPIEntry(applicationName, apiPackage, apiDefinition))
	}

	for _, eventAPIDefinition := range apiPackage.EventDefinitions {
		entries = append(entries, c.toEventServiceEntry(eventAPIDefinition))
	}

	return entries
}

func (c converter) toAPIEntry(applicationName string, apiPackage model.APIPackage, apiDefinition model.APIDefinition) v1alpha1.Entry {

	getApiType := func() string {
		if apiDefinition.APISpec != nil {
			return string(apiDefinition.APISpec.Type)
		}

		return ""
	}
	entry := v1alpha1.Entry{
		ID:                          apiDefinition.ID,
		Name:                        apiDefinition.Name,
		Type:                        SpecAPIType,
		ApiType:                     getApiType(),
		TargetUrl:                   apiDefinition.TargetUrl,
		SpecificationUrl:            "", // Director returns BLOB here
		Credentials:                 c.ToCredentials(applicationName, apiPackage),
		RequestParametersSecretName: "", //TODO: ?
	}

	return entry
}

func (c converter) ToCredentials(applicationName string, apiPackage model.APIPackage) v1alpha1.Credentials {
	result := v1alpha1.Credentials{}
	if apiPackage.DefaultInstanceAuth != nil {
		if apiPackage.DefaultInstanceAuth.Credentials.Oauth != nil {
			result.SecretName = c.credentialsSecretName(applicationName, apiPackage)
			result.Type = secrets.TypeOAuth
			if apiPackage.DefaultInstanceAuth.Credentials.CSRFInfo != nil {
				result.CSRFInfo = &v1alpha1.CSRFInfo{
					TokenEndpointURL: apiPackage.DefaultInstanceAuth.Credentials.CSRFInfo.TokenEndpointURL,
				}
			}
		} else if apiPackage.DefaultInstanceAuth.Credentials.Basic != nil {
			result.SecretName = c.credentialsSecretName(applicationName, apiPackage)
			result.Type = secrets.TypeBasic
		}
	}
	return result
}

func (c converter) credentialsSecretName(applicationName string, apiPackage model.APIPackage) string {
	return fmt.Sprintf("%s-%s", applicationName, apiPackage.ID)
}

func (c converter) toEventServiceEntry(eventsDefinition model.EventAPIDefinition) v1alpha1.Entry {
	return v1alpha1.Entry{
		ID:               eventsDefinition.ID,
		Name:             eventsDefinition.Name,
		Type:             SpecEventsType,
		SpecificationUrl: "", // Director returns BLOB here
	}
}

func (c converter) toCompassMetadata(applicationID string, systemAuthsIDs []string) *v1alpha1.CompassMetadata {
	return &v1alpha1.CompassMetadata{
		ApplicationID: applicationID,
		Authentication: v1alpha1.Authentication{
			ClientIds: systemAuthsIDs,
		},
	}
}
