package config

// ConvertServices converts a set of v1 service configs to compose-go service configs.
//
// This is the main entry point used by BurmillaOS to convert cloud-config service
// definitions (ServiceConfigV1, deserialized via candiedyaml) into compose-go
// ServiceConfig values used by the compose engine.
//
// The conversion is fully backward compatible: all v1 config keys map 1:1 to their
// compose-go equivalents. Invalid services are skipped with a log message rather
// than causing a conversion error, so existing valid configurations are unaffected.
func ConvertServices(v1Services map[string]*ServiceConfigV1) (map[string]*ServiceConfig, error) {
	return ConvertV1Services(v1Services), nil
}
