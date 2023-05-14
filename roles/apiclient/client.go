package apiclient

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/nadavbm/etzba/pkg/debug"
	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/zlog"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Client is an api client
type Client struct {
	Logger *zlog.Logger
	client *http.Client
	// auth contains authentication method and token for api server
	auth *authenticator.ApiAuth
}

// Response is the server response, currently only status and payload
type Response struct {
	Status  string `json:"status"`
	Payload string `json:"payload"`
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

// ExecuteAPIRequest to remote server url by ApiRequest
func (c *Client) ExecuteAPIRequest(req *ApiRequest) (*Response, error) {
	return c.createAPIRequest(req.Url, req.Method, []byte(req.Payload))
}

// createAPIRequest will create GET, POST, PUT or DELETE request for api server url
func (c *Client) createAPIRequest(url, method string, reqBody []byte) (*Response, error) {
	var req *http.Request
	var err error

	switch {
	case method == http.MethodPost:
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create post request")
		}
		req.Header.Set("Content-Type", "application/json")
	case method == http.MethodPut:
		req, err = http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create put request")
		}
		req.Header.Set("Content-Type", "application/json")
	default:
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create get request")
		}
	}

	if c.auth != nil {
		req.Header.Add("Authorization", c.auth.Method+" "+c.auth.Token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		debug.Debug("error1", err)
		return &Response{
			Status: err.Error(),
		}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger.Error("falied to close response", zap.Error(err))
		}
	}()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		debug.Debug("error2", err)
		return &Response{
			Status: err.Error(),
		}, err
	}

	return &Response{
		Status:  resp.Status,
		Payload: string(resBody),
	}, nil
}
