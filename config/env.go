package config

import (
	"os"
)

func GetEnv(key string) (string, error) {
	return os.Getenv(key), nil
}
