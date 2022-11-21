package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
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

func GetProjectID(ctx context.Context) (string, error) {
	credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		return "", fmt.Errorf("google.FindDefaultCredentials failed; %w", err)
	}
	return credentials.ProjectID, nil
}

func GetSecret(ctx context.Context, secretID string) (string, error) {
	clt, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("secretmanager.NewClient failed; %w", err)
	}
	defer clt.Close()
	projectID, err := GetProjectID(ctx)
	if err != nil {
		return "", err
	}
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID),
	}
	resp, err := clt.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("secretmanager.Client.AccessSecretVersion failed; %w", err)
	}
	return string(resp.Payload.Data), nil
}
