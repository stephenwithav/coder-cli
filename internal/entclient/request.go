package entclient

import (
	"bytes"
	"encoding/json"
	"net/http"

	"golang.org/x/xerrors"
)

func (c Client) request(
	method string, path string,
	request interface{},
) (*http.Response, error) {
	client, err := c.http()
	if err != nil {
		return nil, err
	}
	if request == nil {
		request = []byte{}
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, xerrors.Errorf("marshal request: %w", err)
	}
	req, err := http.NewRequest(method, c.BaseURL.String()+path, bytes.NewReader(body))
	if err != nil {
		return nil, xerrors.Errorf("create request: %w", err)
	}
	return client.Do(req)
}

func (c Client) requestBody(
	method string, path string, request interface{}, response interface{},
) error {
	resp, err := c.request(method, path, request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return bodyError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	return nil
}
