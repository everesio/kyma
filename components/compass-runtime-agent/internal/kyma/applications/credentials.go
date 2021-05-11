package applications

import (
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/k8sconsts"
	"github.com/kyma-project/kyma/components/compass-runtime-agent/internal/kyma/model"
)

func ToCredentials(applicationName string, apiPackage model.APIPackage) *Credentials {
	if apiPackage.DefaultInstanceAuth != nil && apiPackage.DefaultInstanceAuth.Credentials != nil {
		if apiPackage.DefaultInstanceAuth.Credentials.Oauth != nil {
			csrfInfo := func(csrfInfo *model.CSRFInfo) *CSRFInfo {
				if csrfInfo != nil {
					return &CSRFInfo{
						TokenEndpointURL: csrfInfo.TokenEndpointURL,
					}
				}
				return nil
			}
			headers := func(params *model.RequestParameters) *map[string][]string {
				if params != nil {
					return params.Headers
				}
				return nil
			}
			queryParameters := func(params *model.RequestParameters) *map[string][]string {
				if params != nil {
					return params.QueryParameters
				}
				return nil
			}
			return &Credentials{
				Type:              CredentialsOAuthType,
				SecretName:        credentialsSecretName(applicationName, apiPackage),
				AuthenticationUrl: apiPackage.DefaultInstanceAuth.Credentials.Oauth.URL,
				CSRFInfo:          csrfInfo(apiPackage.DefaultInstanceAuth.Credentials.CSRFInfo),
				Headers:           headers(apiPackage.DefaultInstanceAuth.Credentials.Oauth.RequestParameters),
				QueryParameters:   queryParameters(apiPackage.DefaultInstanceAuth.Credentials.Oauth.RequestParameters),
			}
		} else if apiPackage.DefaultInstanceAuth.Credentials.Basic != nil {
			return &Credentials{
				Type:       CredentialsBasicType,
				SecretName: credentialsSecretName(applicationName, apiPackage),
			}
		}
	}
	return nil
}

func GetApplicationCredentials(directorApplication model.Application) []Credentials {
	result := make([]Credentials, 0)
	for _, apiPackage := range directorApplication.APIPackages {
		credentials := ToCredentials(directorApplication.Name, apiPackage)
		if credentials != nil {
			result = append(result, *credentials)
		}
	}
	return result
}

func credentialsSecretName(applicationName string, apiPackage model.APIPackage) string {
	return k8sconsts.GetResourceName(applicationName, apiPackage.ID)
}
