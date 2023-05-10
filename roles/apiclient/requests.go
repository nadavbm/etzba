package apiclient

type ApiRequest struct {
	Url             string          `json:"url,omitempty"`
	Method          string          `json:"method,omitempty"`
	Payload         string          `json:"payload,omitempty"`
	EndpointFile    string          `json:"endpointFile,omitempty"`
	EndpointPattern EndpointPattern `json:"endpointPattern,omitempty"`
	Weight          int             `json:"weight,omitempty"`
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
