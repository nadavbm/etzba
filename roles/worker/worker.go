package worker

import (
	"time"

	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/sqlclient"
	"github.com/nadavbm/zlog"
	"github.com/pkg/errors"
)

// APIWorker will get an assignment from the Scheduler
type APIWorker struct {
	Logger    *zlog.Logger
	ApiClient *apiclient.Client
}

// NewWorker creates an instance of a worker
func NewAPIWorker(logger *zlog.Logger, secretFile string) (*APIWorker, error) {
	apiClient, err := apiclient.NewClient(logger, secretFile)
	if err != nil {
		return nil, err
	}
	return &APIWorker{
		Logger:    logger,
		ApiClient: apiClient,
	}, nil
}

func (w *APIWorker) GetAPIRequestDuration(assignment *Assignment) (time.Duration, error) {
	// start to count api request duration
	start := time.Now()

	_, err := w.ApiClient.CreateAPIRequest()
	if err != nil {
		return time.Since(start), err
	}

	return time.Since(start), nil
}

// SQLWorker will get an assignment from the Scheduler
type SQLWorker struct {
	Logger    *zlog.Logger
	SqlClient *sqlclient.Client
}

// NewSQLWorker creates an instance of a worker
func NewSQLWorker(logger *zlog.Logger, secretFile string) (*SQLWorker, error) {
	sqlClient, err := sqlclient.NewClient(logger, secretFile)
	if err != nil {
		return nil, err
	}
	return &SQLWorker{
		Logger:    logger,
		SqlClient: sqlClient,
	}, nil
}

func (w *SQLWorker) GetSQLQueryDuration(assignment *Assignment) (time.Duration, error) {
	// start to count sql query duration
	start := time.Now()
	if err := w.SqlClient.ExecuteQuery(translateAssignmentToQueryBuilder(assignment)); err != nil {
		return time.Since(start), errors.Wrap(err, "worker could not execute query")
	}

	return time.Since(start), nil
}
