package manifest

type ManifestInformation struct {
	Version string `json:"version"`
	BuildNumber string `json:"build"`
}

var information ManifestInformation

func SetApplicationInfo(version string, buildId string) {
	information.Version = version
	information.BuildNumber = buildId
}

// Get application version information
func GetApplicationInfo() *ManifestInformation {
	return &information
}