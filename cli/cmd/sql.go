package cmd

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nadavbm/etzba/pkg/debug"
	"github.com/nadavbm/etzba/pkg/printer"
	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/etzba/roles/scheduler"
	"github.com/nadavbm/zlog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	sqlCmd = &cobra.Command{
		Use:       "sql",
		Short:     "Start benchmarking your sql instance",
		Long:      `Start benchmarking db, defining the number of workers and csv file input`,
		ValidArgs: validArgs,
		Run:       benchmarkSql,
	}
)

func benchmarkSql(cmd *cobra.Command, args []string) {
	logger := zlog.New()

	jobDuration, err := setDurationFromString(duration)
	if err != nil {
		logger.Fatal("could set job duration")
	}

	s, err := scheduler.NewScheduler(logger, jobDuration, "sql", configFile, helpersFile, rps, workersCount, Verbose)
	if err != nil {
		logger.Fatal("could not create a scheduler instance")
	}

	auth, err := s.Authenticator.GetSQLAuth()
	if err != nil {
		logger.Fatal("could not get secret for sql auth from config file")
	}

	debug.Debug("start connection pool")
	pool, err := pgxpool.Connect(context.Background(), getConnectionString(auth))
	if err != nil {
		logger.Fatal("could not create a db connection pool")
	}
	defer pool.Close()

	s.ConnectionPool = pool
	debug.Debug("connection pool config", s.ConnectionPool.Config().ConnConfig)

	var result *scheduler.Result
	if duration != "" {
		result, err = s.ExecuteJobByDuration()
		if err != nil {
			s.Logger.Fatal("could not start execution", zap.Error(err))
		}
	} else {
		if result, err = s.ExecuteJobUntilCompletion(); err != nil {
			s.Logger.Fatal("could not start execution")
		}
	}

	printer.PrintToTerminal(result, false)
}

//
// ---------------------------------------------------------------------------------- helpers -----------------------------------------------------------------------------------------------------------
//

// getConnectionString return a connection string based on environment vars
func getConnectionString(auth *authenticator.SqlAuth) string {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", auth.User, auth.Password, auth.Host, auth.Port, auth.Database)
	return conn
}
