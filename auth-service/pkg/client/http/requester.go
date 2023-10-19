package client

import (
	"encoding/json"
	"net/http"
	"time"
)

var defaultTimeout = 10 * time.Second

type IDo interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type IRequester interface {
	Get(url string, dst any) error
	// Post(url string, dst any, payload io.Reader) error
	// PUT
	// PATCH
}

type AuthClient struct {
	c IDo
}

func NewClient() *AuthClient {
	var _ IDo = (*http.Client)(nil)
	c := &http.Client{
		Timeout: defaultTimeout,
	}
	return &AuthClient{c: c}
}

func (authClient AuthClient) Get(url string, dst any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err // TODO: wrap and add context
	}

	httpResp, err := authClient.c.Do(req)
	if err != nil {
		return err // TODO: wrap and add context
	}

	defer httpResp.Body.Close()
	err = json.NewDecoder(httpResp.Body).Decode(&dst)
	if err != nil {
		return err // TODO: wrap and add context
	}

	return nil
}
