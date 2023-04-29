package apiclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/zlog"
	"go.uber.org/zap"
)

type Client struct {
	Logger      *zlog.Logger
	client      *http.Client
	auth        *authenticator.ApiAuth
	requestFile string
}

// NewClient creates an instance of api client
func NewClient(logger *zlog.Logger, secretFile string) (*Client, error) {
	a := authenticator.NewAuthenticator(logger, secretFile)
	auth, err := a.GetAPIAuth()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return &Client{
		Logger: logger,
		client: client,
		auth:   auth,
	}, nil
}

func (c *Client) CreateAPIRequest() ([]byte, error) {
	encodedConfig, err := c.readApiRequestSpecificationsFromJsonFile()
	if err != nil {
		return nil, err
	}

	var req ApiRequest
	if err := json.Unmarshal(encodedConfig, &req); err != nil {
		return nil, err
	}

	switch {
	case req.Method == http.MethodPost:
		res, err := c.postRequest(req.Payload, req.Url)
		if err != nil {
			return nil, err
		}
		return res, nil
	case req.Method == http.MethodPut:
		res, err := c.putRequest(req.Payload, req.Url)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		res, err := c.getRequest(req.Url)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (c *Client) getRequest(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger.Error("falied to close response", zap.Error(err))
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) postRequest(encodedBody []byte, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(encodedBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "applications/json")
	req.Header.Add("Authorization", c.auth.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger.Error("falied to close response", zap.Error(err))
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) putRequest(encodedBody []byte, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(encodedBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "applications/json")
	req.Header.Add("Authorization", c.auth.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger.Error("falied to close response", zap.Error(err))
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
