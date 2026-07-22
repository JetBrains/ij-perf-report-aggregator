package util

import (
	"os"
	"strconv"
)

// empty values are treated as unset, matching the previous go.deanishe.net/env behavior

func GetEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func GetEnvAsInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func GetEnvAsBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func GetEnvOrFile(envName string, file string) (string, error) {
	v := os.Getenv(envName)
	if v == "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}
		return string(b), err
	}
	return v, nil
}

func GetEnvOrFileOrPanic(envName string, file string) string {
	v := os.Getenv(envName)
	if v == "" {
		b, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		return string(b)
	}
	return v
}
