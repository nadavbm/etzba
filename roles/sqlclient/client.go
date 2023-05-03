package sqlclient

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/nadavbm/etzba/pkg/debug"
	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/zlog"
)

var conn *pgx.Conn

type Client struct {
	Logger *zlog.Logger
	auth   *authenticator.SqlAuth
}

// NewClient creates an instance of sql client
func NewClient(logger *zlog.Logger, secretFile string) (*Client, error) {
	a := authenticator.NewAuthenticator(logger, secretFile)
	auth, err := a.GetSQLAuth()
	if err != nil {
		return nil, err
	}

	return &Client{
		Logger: logger,
		auth:   auth,
	}, nil
}

func (c *Client) ExecuteQuery(b *QueryBuilder) error {
	query := toSQL(b)
	debug.Debug("query", query)
	switch {
	case b.Command == "INSERT" || b.Command == "insert":
		return c.execQuery(query)
	case b.Command == "UPDATE" || b.Command == "update":
		return c.execQuery(query)
	case b.Command == "DELETE" || b.Command == "delete":
		return c.execQuery(query)
	default:
		return c.selectQuery(query)
	}

}

// getConnectionString return a connection string based on environment vars
func getConnectionString(auth *authenticator.SqlAuth) string {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", auth.User, auth.Password, auth.Host, auth.Port, auth.Database)
	return conn
}
