package utilgcp

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func ProjectID(ctx context.Context) (string, error) {
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
	projectID, err := ProjectID(ctx)
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
