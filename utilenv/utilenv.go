package utilenv

import (
	"fmt"
	"os"
	"strconv"
)

func RequiredVar(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("you must define %s env var", name)
	}
	return val, nil
}

func GetPort() (int, error) {
	portStr, err := RequiredVar("PORT")
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi failed; %s; %w", portStr, err)
	}
	return port, nil
}
