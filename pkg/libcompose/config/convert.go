package config

// ConvertServices converts a set of v1 service configs to compose-go service configs
func ConvertServices(v1Services map[string]*ServiceConfigV1) (map[string]*ServiceConfig, error) {
	return ConvertV1Services(v1Services), nil
}
