package config

import "os"

func unsetenv(key string) error {
	return os.Unsetenv(key)
}
