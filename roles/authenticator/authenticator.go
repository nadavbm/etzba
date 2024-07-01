package authenticator

import (
	"encoding/json"
	"strings"

	"github.com/nadavbm/etzba/pkg/env"
	"github.com/nadavbm/etzba/pkg/filer"
	"github.com/nadavbm/zlog"
	"gopkg.in/yaml.v3"
)

// Authenticator takes a secret file and authenticate to sql or api server
type Authenticator struct {
	Logger     *zlog.Logger
	reader     *filer.Reader
	SecretFile string
}

// NewAuthenticator creates an instance of authenticator
func NewAuthenticator(logger *zlog.Logger, secretFile string) *Authenticator {
	return &Authenticator{
		Logger:     logger,
		reader:     filer.NewReader(logger),
		SecretFile: secretFile,
	}
}

// Secret for sql or api authentication is taken from "--secret=file.json" arg of the command line
type Secret struct {
	ApiAuth ApiAuth `json:"apiAuth,omitempty" yaml:"apiAuth"`
	SqlAuth SqlAuth `json:"sqlAuth,omitempty" yaml:"sqlAuth"`
}

// ApiAuth is an api server authentication
type ApiAuth struct {
	// Method is the authentication method, e.g. Bearer or ApiKey
	Method string `json:"method,omitempty" yaml:"method,omitempty"`
	// Token is the authentication token (Bearer token or API key value)
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

// SqlAuth is a sql server authentication params
type SqlAuth struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	Database string `json:"database,omitempty" yaml:"database,omitempty"`
	User     string `json:"user,omitempty" yaml:"user,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

//
// ------------------------------------------------------------------------------------------ sql auth -------------------------------------------------------------------------------------------------
//

// GetSQLAuth returns sql authentication params from a secret
func (a *Authenticator) GetSQLAuth() (*SqlAuth, error) {
	if getSQLAuthFromEnv() != nil {
		return getSQLAuthFromEnv(), nil
	}
	return a.getSQLAuthFromFile()
}

// getSQLAuthFromEnv gets secrets from environment variable (if set)
func getSQLAuthFromEnv() *SqlAuth {
	var auth SqlAuth
	if env.DatabaseDB != "" && env.DatabaseHost != "" && env.DatabasePort != 0 && env.DatabaseUser != "" && env.DatabasePass != "" {
		auth.Database = env.DatabaseDB
		auth.Host = env.DatabaseHost
		auth.Port = env.DatabasePort
		auth.User = env.DatabaseUser
		auth.Password = env.DatabasePass
		return &auth
	}

	return nil
}

// getSQLAuthFromFile parse secret from a file and returns sql auth
func (a *Authenticator) getSQLAuthFromFile() (*SqlAuth, error) {
	secret, err := a.parseSecret()
	if err != nil {
		return nil, err
	}
	return &secret.SqlAuth, nil
}

//
// ------------------------------------------------------------------------------------------ api auth -------------------------------------------------------------------------------------------------
//

// GetAPIAuth returns api authentication params from a secret
func (a *Authenticator) GetAPIAuth() (*ApiAuth, error) {
	auth := getAPIAuthFromEnv()
	if auth != nil {
		return auth, nil
	}
	return a.getAPIAuthFromFile()
}

// getAPIAuthFromEnv gets secrets from environment variable (if set)
func getAPIAuthFromEnv() *ApiAuth {
	var auth ApiAuth
	if env.ApiToken != "" && env.ApiAuthMethod != "" {
		auth.Method = env.ApiAuthMethod
		auth.Token = env.ApiToken
		return &auth
	}
	return nil
}

// GetAPIAuth returns api authentication params from a secret
func (a *Authenticator) getAPIAuthFromFile() (*ApiAuth, error) {
	secret, err := a.parseSecret()
	if err != nil {
		return nil, err
	}
	return &secret.ApiAuth, nil
}

// parseSecret create a secret from json file
func (a *Authenticator) parseSecret() (*Secret, error) {
	bs, err := a.reader.ReadFile(a.SecretFile)
	if err != nil {
		return nil, err
	}

	var s Secret
	switch {
	case strings.HasSuffix(a.SecretFile, ".json"):
		if err := json.Unmarshal(bs, &s); err != nil {
			return nil, err
		}
	case strings.HasSuffix(a.SecretFile, ".yaml"):
		if err := yaml.Unmarshal(bs, &s); err != nil {
			return nil, err
		}
	}

	return &s, nil
}
