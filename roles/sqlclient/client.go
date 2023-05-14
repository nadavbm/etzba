package sqlclient

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/zlog"
)

var conn *pgx.Conn

// Client is an sql client that can execute db queries
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

// ExecuteQuery gets an sql query from builder (after translate assignment) and execute the query
func (c *Client) ExecuteQuery(b *QueryBuilder) error {
	query := ToSQL(b)
	switch {
	case b.Command == "INSERT" || b.Command == "insert":
		return c.executeQuery(query)
	case b.Command == "UPDATE" || b.Command == "update":
		return c.executeQuery(query)
	case b.Command == "DELETE" || b.Command == "delete":
		return c.executeQuery(query)
	default:
		return c.executeSelectQuery(query)
	}

}

// getConnectionString return a connection string based on environment vars
func getConnectionString(auth *authenticator.SqlAuth) string {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", auth.User, auth.Password, auth.Host, auth.Port, auth.Database)
	return conn
}
