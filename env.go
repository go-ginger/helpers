package helpers

import "os"

//GetEnv is used in config file for get environment var
//from os and fallback to default value if it is empty.
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
