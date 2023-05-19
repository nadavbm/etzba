package sqlclient

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	command := strings.ToUpper(querySpec.Command)
	constraint := ""
	if querySpec.Constraint != "" {
		constraint = fmt.Sprintf("WHERE %s", querySpec.Constraint)
	}

	switch {
	case command == "INSERT":
		coulmnsRef := strings.Split(querySpec.ColumnsRef, " ")
		values := strings.Split(querySpec.Values, " ")
		columns := ""
		vals := ""
		for i := 0; i < len(coulmnsRef); i++ {
			columns += fmt.Sprintf("%s,", coulmnsRef[i])
			vals += fmt.Sprintf("%s,", values[i])
		}

		return fmt.Sprintf("%s INTO %s (%s) VALUES (%s)", command, querySpec.Table, columns, vals)
	case command == "UPDATE":
		coulmnsRef := strings.Split(querySpec.ColumnsRef, " ")
		values := strings.Split(querySpec.Values, " ")
		setColumnsAndValues := ""
		for i := 0; i < len(coulmnsRef); i++ {
			setColumnsAndValues += fmt.Sprintf("%s = %s, ", coulmnsRef[i], values[i])
		}

		return fmt.Sprintf("%s %s SET %s %s", command, querySpec.Table, setColumnsAndValues, constraint)
	case command == "DELETE":
		return fmt.Sprintf("%s FROM %s %s", command, querySpec.Table, constraint)
	default:
		return fmt.Sprintf("%s * FROM %s %s", command, querySpec.Table, constraint)
	}
}
