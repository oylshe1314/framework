package util

import "os"

func ReadFileString(name string) (string, error) {
	buf, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
