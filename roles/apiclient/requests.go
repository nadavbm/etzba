package apiclient

import (
	"io/ioutil"
	"os"

	"go.uber.org/zap"
)

type ApiRequest struct {
	Url             string          `json:"url,omitempty"`
	Method          string          `json:"method,omitempty"`
	Payload         []byte          `json:"payload,omitempty"`
	EndpointFile    string          `json:"endpointFile,omitempty"`
	EndpointPattern EndpointPattern `json:"endpointPattern,omitempty"`
	Weight          int             `json:"weight,omitempty"`
}

type EndpointPattern struct {
	Length     int    `json:"length,omitempty"`
	Occurences int    `json:"occurences,omitempty"`
	Regex      string `json:"regex,omitempty"`
}

func (c *Client) readApiRequestSpecificationsFromJsonFile() ([]byte, error) {
	jsonFile, err := os.Open(c.requestFile)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := jsonFile.Close(); err != nil {
			c.Logger.Error("failed to close json file", zap.Error(err))
		}
	}()

	return ioutil.ReadAll(jsonFile)
}

type ApiResponse struct {
	Status  int    `json:"status,omitempty"`
	Payload []byte `json:"payload,omitempty"`
}
