package sqlclient

import (
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/zlog"
)

// Client is an sql client that can execute db queries
type Client struct {
	Logger *zlog.Logger
}

// NewClient creates an instance of sql client
func NewClient(logger *zlog.Logger) (*Client, error) {
	return &Client{
		Logger: logger,
	}, nil
}

// ExecuteQuery gets an sql query from builder (after translate assignment) and execute the query
func (c *Client) ExecuteQuery(b *QueryBuilder, conn *pgxpool.Conn) error {
	query := ToSQL(b)
	switch {
	case b.Command == "INSERT" || b.Command == "insert":
		return c.executeQuery(query, conn)
	case b.Command == "UPDATE" || b.Command == "update":
		return c.executeQuery(query, conn)
	case b.Command == "DELETE" || b.Command == "delete":
		return c.executeQuery(query, conn)
	default:
		return c.executeSelectQuery(query, conn)
	}

}

// GetConnectionString return a connection string based on environment vars
func GetConnectionString(auth *authenticator.SqlAuth) string {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", auth.User, auth.Password, auth.Host, auth.Port, auth.Database)
	return conn
}
