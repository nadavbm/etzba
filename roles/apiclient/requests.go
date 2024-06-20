package apiclient

import (
	"net/http"
)

type ApiRequest struct {
	Url             string          `json:"url,omitempty" yaml:"url"`
	Method          string          `json:"method,omitempty" yaml:"method"`
	Payload         string          `json:"payload,omitempty" yaml:"payload"`
	EndpointFile    string          `json:"endpointFile,omitempty" yaml:"endpointFile"`
	EndpointPattern EndpointPattern `json:"endpointPattern,omitempty" yaml:"endpointPattern"`
	Weight          int             `json:"weight,omitempty" yaml:"weight"`
}

type EndpointPattern struct {
	Length     int    `json:"length,omitempty"`
	Occurences int    `json:"occurences,omitempty"`
	Regex      string `json:"regex,omitempty"`
}

// Response is the server response, currently only status and payload
type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Headers http.Header `json:"headers"`
	Payload string      `json:"payload"`
}
