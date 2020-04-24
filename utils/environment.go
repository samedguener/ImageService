package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// Environment ...
	Environment EnvironmentVariable

	// GoogleApplicationCredentials ...
	GoogleApplicationCredentials EnvironmentVariable

	// AuthenticationMethod ...
	AuthenticationMethod EnvironmentVariable

	// BucketName ...
	BucketName EnvironmentVariable

	// GCPProjectID ...
	GCPProjectID EnvironmentVariable

	// ImageUploadToGCPTimeout ...
	ImageUploadToGCPTimeout EnvironmentVariable

	// ImageAccessEndpoint ...
	ImageAccessEndpoint EnvironmentVariable
)

// EnvironmentVariable ...
type EnvironmentVariable struct {
	Name  string
	Value string
	IsSet bool
}

// InitEnvironmentVariables ...
func InitEnvironmentVariables() {
	Environment = GetEnvironmentVariable("ENVIRONMENT", true, "")

	BucketName = GetEnvironmentVariable("BUCKET_NAME", true, "")

	GCPProjectID = GetEnvironmentVariable("PROJECT_ID", true, "")

	ImageAccessEndpoint = GetEnvironmentVariable("IMAGE_ACCESS_ENDPOINT", true, "")

	AuthenticationMethod = GetEnvironmentVariable("AUTHENTICATION_METHOD", false, "firebase")

	ImageUploadToGCPTimeout = GetEnvironmentVariable("TIMEOUT_IMAGE_UPLOAD_GCP", false, "2s")

	if Environment.Value == "dev" {
		// using Google ADC in other landscapes
		GoogleApplicationCredentials = GetEnvironmentVariable("GOOGLE_APPLICATION_CREDENTIALS", false, "./service_account.json")
	}
}

// GetEnvironmentVariable ...
func GetEnvironmentVariable(name string, required bool, defaultValue string) EnvironmentVariable {

	value, isSet := os.LookupEnv(name)

	if required && (value == "" || isSet == false) {
		logrus.Errorf("error: %s environment variable not set.\n", name)
		os.Exit(1)
	} else if value == "" || isSet == false {
		logrus.Warnf("warn: %s environment is not set defaulting to %s", name, defaultValue)
		value = defaultValue
	}

	return EnvironmentVariable{Name: name, Value: value, IsSet: isSet}
}
