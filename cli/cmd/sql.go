package cmd

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nadavbm/etzba/pkg/printer"
	"github.com/nadavbm/etzba/roles/authenticator"
	"github.com/nadavbm/etzba/roles/common"
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
		logger.Fatal("could set job duration", zap.Error(err))
	}

	settings := common.GetSettings(jobDuration, "sql", configFile, helpersFile, rps, workersCount, Verbose)
	s, err := scheduler.NewScheduler(logger, settings)
	if err != nil {
		logger.Fatal("could not create a scheduler instance", zap.Error(err))
	}

	auth, err := s.Authenticator.GetSQLAuth()
	if err != nil {
		logger.Fatal("could not get secret for sql auth from config file", zap.Error(err))
	}

	pool, err := pgxpool.Connect(context.Background(), getConnectionString(auth))
	if err != nil {
		logger.Fatal("could not create a db connection pool", zap.Error(err))
	}
	defer pool.Close()

	s.ConnectionPool = pool

	var result *scheduler.Result
	if duration != "" {
		result, err = s.ExecuteJobByDuration()
		if err != nil {
			s.Logger.Fatal("could not execute job by duration", zap.Error(err))
		}
	} else {
		if result, err = s.ExecuteJobUntilCompletion(); err != nil {
			s.Logger.Fatal("could not execute job until completion", zap.Error(err))
		}
	}

	printer.PrintToTerminal(result, false)
}

//
// ---------------------------------------------------------------------------------- helpers -----------------------------------------------------------------------------------------------------------
//

func setDBConnectionPool(auth string, workersCount int32) (*pgxpool.Pool, error) {
	// TODO: using parse config create errors and issue with cli.
	// change rps to qps for sql and work on accuracy of queries per second
	connConf, err := pgxpool.ParseConfig(auth)
	connConf.MaxConns = workersCount
	connConf.MinConns = 10

	pool, err := pgxpool.ConnectConfig(context.TODO(), connConf)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// getConnectionString return a connection string based on environment vars
func getConnectionString(auth *authenticator.SqlAuth) string {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", auth.User, auth.Password, auth.Host, auth.Port, auth.Database)
	return conn
}
