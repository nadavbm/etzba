package apiclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/zlog"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Client struct {
	Logger      *zlog.Logger
	client      *http.Client
	auth        *authenticator.ApiAuth
	requestFile string
}

type Response struct {
	Status   int    `json:"status"`
	Payload  string `json:"payload"`
	APIError string
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

func (c *Client) ExecuteAPIRequest() (*Response, error) {
	encodedConfig, err := c.readApiRequestSpecificationsFromJsonFile()
	if err != nil {
		return nil, err
	}

	var req ApiRequest
	if err := json.Unmarshal(encodedConfig, &req); err != nil {
		return nil, err
	}

	return c.createAPIRequest(req.Url, req.Method)
}

func (c *Client) createAPIRequest(url, method string) (*Response, error) {
	var req *http.Request
	var err error
	switch {
	case req.Method == http.MethodPost:
		req, err = http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create post request")
		}
	case req.Method == http.MethodPut:
		req, err = http.NewRequest(http.MethodPut, url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create put request")
		}
	default:
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create get request")
		}
	}

	if c.auth != nil {
		req.Header.Add("Authorization", c.auth.AuthMethod+" "+c.auth.Token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return &Response{
			Status:   resp.StatusCode,
			APIError: err.Error(),
		}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger.Error("falied to close response", zap.Error(err))
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			Status:   resp.StatusCode,
			APIError: err.Error(),
		}, err
	}

	return &Response{
		Status:  resp.StatusCode,
		Payload: string(body),
	}, nil
}
