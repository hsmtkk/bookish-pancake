package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func RequiredEnvVar(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("you must define %s env var", name)
	}
	return val, nil
}

func GetPort() (int, error) {
	portStr, err := RequiredEnvVar("PORT")
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi failed; %s; %w", portStr, err)
	}
	return port, nil
}

func GetHostFromURL(urlStr string) (string, error) {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("url.Parse failed; %w", err)
	}
	return parsed.Host, nil
}

func HandleHTTPResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 HTTP status code; %d; %s", resp.StatusCode, resp.Status)
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed; %w", err)
	}
	return bs, nil
}
