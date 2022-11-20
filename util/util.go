package util

import (
	"fmt"
	"io"
	"net/http"
)

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
