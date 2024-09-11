package util

import "os"

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
