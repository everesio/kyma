package k8sconsts

import (
	"fmt"
)

const (
	resourceNamePrefixFormat = "%s-"
	maxResourceNameLength    = 63 // Kubernetes limit for services
	uuidLength               = 36 // UUID has 36 characters
)

// GetResourceName returns resource name with given ID
func GetResourceName(application, id string) string {
	return getResourceNamePrefix(application) + id
}

func getResourceNamePrefix(application string) string {
	truncatedApplicaton := truncateApplication(application)
	return fmt.Sprintf(resourceNamePrefixFormat, truncatedApplicaton)
}

func truncateApplication(application string) string {
	maxResourceNamePrefixLength := maxResourceNameLength - uuidLength
	testResourceNamePrefix := fmt.Sprintf(resourceNamePrefixFormat, application)
	testResourceNamePrefixLength := len(testResourceNamePrefix)

	overflowLength := testResourceNamePrefixLength - maxResourceNamePrefixLength

	if overflowLength > 0 {
		newApplicationLength := len(application) - overflowLength
		return application[0:newApplicationLength]
	}
	return application
}
