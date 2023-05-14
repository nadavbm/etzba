package authenticator

import (
	"encoding/json"

	"github.com/nadavbm/etzba/pkg/reader"
	"github.com/nadavbm/zlog"
)

// Authenticator takes a secret file and authenticate to sql or api server
type Authenticator struct {
	Logger     *zlog.Logger
	reader     *reader.Reader
	SecretFile string
}

// NewAuthenticator creates an instance of authenticator
func NewAuthenticator(logger *zlog.Logger, secretFile string) *Authenticator {
	return &Authenticator{
		Logger:     logger,
		reader:     reader.NewReader(logger),
		SecretFile: secretFile,
	}
}

// Secret for sql or api authentication is taken from "--secret=file.json" arg of the command line
type Secret struct {
	ApiAuth ApiAuth `json:"apiAuth,omitempty"`
	SqlAuth SqlAuth `json:"sqlAuth,omitempty"`
}

// ApiAuth is an api server authentication
type ApiAuth struct {
	// Method is the authentication method, e.g. Bearer or ApiKey
	Method string `json:"method,omitempty"`
	// Toekn is the authentication token (Bearer token or API key value)
	Token string `json:"token,omitempty"`
}

// SqlAuth is a sql server authentication params
type SqlAuth struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

// GetSQLAuth returns sql authentication params from a secret
func (a *Authenticator) GetSQLAuth() (*SqlAuth, error) {
	secret, err := a.parseSecret()
	if err != nil {
		return nil, err
	}

	return &secret.SqlAuth, nil
}

// GetAPIAuth returns api authentication params from a secret
func (a *Authenticator) GetAPIAuth() (*ApiAuth, error) {
	secret, err := a.parseSecret()
	if err != nil {
		return nil, err
	}
	return &secret.ApiAuth, nil
}

// parseSecret create a secret from json file
func (a *Authenticator) parseSecret() (*Secret, error) {
	bs, err := a.reader.ReadJSONFile(a.SecretFile)
	if err != nil {
		return nil, err
	}

	var s Secret
	if err := json.Unmarshal(bs, &s); err != nil {
		return nil, err
	}

	return &s, nil
}
