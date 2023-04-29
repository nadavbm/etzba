package sqlclient

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/zlog"
	"go.uber.org/zap"
)

var conn *pgx.Conn

type Client struct {
	Logger      *zlog.Logger
	auth        *authenticator.SqlAuth
	QueriesFile string
}

// NewClient creates an instance of sql client
func NewClient(logger *zlog.Logger, secretFile, queriesFile string) (*Client, error) {
	a := authenticator.NewAuthenticator(logger, secretFile)
	auth, err := a.GetSQLAuth()
	if err != nil {
		return nil, err
	}

	return &Client{
		Logger:      logger,
		auth:        auth,
		QueriesFile: queriesFile,
	}, nil
}

func (c *Client) ExecuteQueries() error {
	data, err := c.readCSVFile()
	if err != nil {
		return err
	}

	builders, err := parseQuerySpecifications(data)
	if err != nil {
		return err
	}

	for _, b := range builders {
		query := toSQL(b)
		switch {
		case b.Command == "INSERT" || b.Command == "insert":
			return c.execQuery(query)
		case b.Command == "UPDATE" || b.Command == "update":
			return c.execQuery(query)
		case b.Command == "DELETE" || b.Command == "delete":
			return c.execQuery(query)
		default:
			if err := c.selectQuery(query); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

// getConnectionString return a connection string based on environment vars
func getConnectionString(auth *authenticator.SqlAuth) string {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", auth.User, auth.Password, auth.Host, auth.Port, auth.Database)
	return conn
}

//
// ---------------------------------------------------------------------------------------- csv reader ------------------------------------------------------------------------------
//

// parseQuerySpecifications from csv data will return query specification to build an sql query
func parseQuerySpecifications(data [][]string) ([]QueryBuilder, error) {
	var builders []QueryBuilder
	for i, row := range data {
		var b QueryBuilder
		if i > 0 {
			b.Command = row[0]
			b.Table = row[1]
			b.Constraint = row[2]
			b.ColumnsRef = row[3]
			b.Values = row[4]

			builders = append(builders, b)
		}
	}

	return builders, nil
}

// ReadCSVFile will get a csv file path and use csv reader
func (c *Client) readCSVFile() ([][]string, error) {
	f, err := os.Open(c.QueriesFile)
	if err != nil {
		c.Logger.Fatal("failed to open csv file", zap.Error(err))
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	return csvReader.ReadAll()
}
