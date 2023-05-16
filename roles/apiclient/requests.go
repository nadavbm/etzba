package apiclient

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

type ApiResponse struct {
	Status  int    `json:"status,omitempty"`
	Payload []byte `json:"payload,omitempty"`
}
