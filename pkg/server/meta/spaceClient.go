package meta

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type SpacePackagesClient struct {
	spaceUrl   string
	spaceToken string
}

func NewSpacePackagesClient(spaceUrl, spaceToken string) *SpacePackagesClient {
	return &SpacePackagesClient{
		spaceUrl:   spaceUrl,
		spaceToken: spaceToken,
	}
}

func (client *SpacePackagesClient) UploadFile(ctx context.Context, project string, packageName string, remoteFolder string, fileName string, file []byte) error {
	endpoint := fmt.Sprintf("/files/p/%s/%s/%s/%s", project, packageName, remoteFolder, fileName)

	_, err := client.doRequest(ctx, endpoint, http.MethodPut, bytes.NewReader(file), nil)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	return nil
}

func (client *SpacePackagesClient) doRequest(ctx context.Context, endpoint string, method string, body io.Reader, headers map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", client.spaceUrl, endpoint)
	slog.Info("Space request", "url", url, "method", method)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.spaceToken)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("request failed with status: %s. Body: %s", resp.Status, bodyString)
	}

	return bodyBytes, nil
}
