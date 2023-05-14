package sqlclient

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

// executeSelectQuery run SELECT by query specifications
func (c *Client) executeSelectQuery(querySpecs string) error {
	ctx := context.TODO()
	conn, err := pgx.Connect(ctx, getConnectionString(c.auth))
	if err != nil {
		c.Logger.Fatal("could not connet to db")
		return errors.Wrap(err, "could not connect to database")
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(context.TODO(), querySpecs)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return errors.Wrap(err, "could not run query")
	}
	defer rows.Close()

	return nil
}

// execQuery execute queries of INSERT, UPDATE and DELETE
func (c *Client) executeQuery(querySpecs string) error {
	_, err := conn.Exec(context.Background(), querySpecs)
	return err
}

//
// ------------------------------------------------------------------------- query builder -------------------------------------------------------
//

// QueryBuilder will transform query specs from csv file into sql query according to the command provided in the csv file
type QueryBuilder struct {
	Command    string `csv:"command"`
	Table      string `csv:"table"`
	Constraint string `csv:"constraint"`
	ColumnsRef string `csv:"colomnref"`
	Values     string `csv:"values"`
}

// ToSQL get a query builder and return an sql query
func ToSQL(querySpec *QueryBuilder) string {
	switch {
	case querySpec.Command == "INSERT":
		return fmt.Sprintf("%s INTO %s (%s) VALUES (%s)", querySpec.Command, querySpec.Table, querySpec.ColumnsRef, querySpec.Values)
	default:
		return fmt.Sprintf("%s * FROM %s WHERE %s", querySpec.Command, querySpec.Table, querySpec.Constraint)
	}
}
