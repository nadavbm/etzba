package authenticator

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/nadavbm/zlog"
	"go.uber.org/zap"
)

type Authenticator struct {
	Logger     *zlog.Logger
	SecretFile string
}

func NewAuthenticator(logger *zlog.Logger, secretFile string) *Authenticator {
	return &Authenticator{
		Logger:     logger,
		SecretFile: secretFile,
	}
}

type Secret struct {
	ApiAuth ApiAuth `json:"apiAuth,omitempty"`
	SqlAuth SqlAuth `json:"sqlAuth,omitempty"`
}

type ApiAuth struct {
	ApiKey string `json:"apiKey,omitempty"`
	Token  string `json:"token,omitempty"`
}

type SqlAuth struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

func (a *Authenticator) GetSQLAuth() (*SqlAuth, error) {
	secret, err := a.parseSecret()
	if err != nil {
		return nil, err
	}

	return &secret.SqlAuth, nil
}

func (a *Authenticator) GetAPIAuth() (*ApiAuth, error) {
	secret, err := a.parseSecret()
	if err != nil {
		return nil, err
	}
	return &secret.ApiAuth, nil
}

func (a *Authenticator) parseSecret() (*Secret, error) {
	bs, err := a.readSecretFile()
	if err != nil {
		return nil, err
	}

	var s Secret
	if err := json.Unmarshal(bs, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (a *Authenticator) readSecretFile() ([]byte, error) {
	jsonFile, err := os.Open(a.SecretFile)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := jsonFile.Close(); err != nil {
			a.Logger.Error("failed to close json file", zap.Error(err))
		}
	}()

	return ioutil.ReadAll(jsonFile)
}
