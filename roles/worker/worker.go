package worker

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nadavbm/etzba/roles/apiclient"
	"github.com/nadavbm/etzba/roles/prompusher"
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

// GetAPIRequestDuration will execute api request and measure the duration since the request started until response
func (w *APIWorker) GetAPIRequestDuration(assignment *Assignment) (time.Duration, *apiclient.Response) {
	// start to count api request duration
	start := time.Now()

	response, err := w.ApiClient.ExecuteAPIRequest(translateAssignmentToAPIRequest(assignment))
	if err != nil {
		return time.Since(start), nil
	}

	pushRequestDurationToPrometheus([]string{"api"}, float64(time.Since(start)))
	return time.Since(start), response
}

// SQLWorker will get an assignment from the Scheduler
type SQLWorker struct {
	Logger         *zlog.Logger
	SqlClient      *sqlclient.Client
	ConnectionPool *pgxpool.Pool
}

// NewSQLWorker creates an instance of a worker
func NewSQLWorker(logger *zlog.Logger, secretFile string, pool *pgxpool.Pool) (*SQLWorker, error) {
	sqlClient, err := sqlclient.NewClient(logger)
	if err != nil {
		return nil, err
	}

	return &SQLWorker{
		Logger:         logger,
		SqlClient:      sqlClient,
		ConnectionPool: pool,
	}, nil
}

// GetSQLQueryDuration will execute a query in the database and measure the duration it takes
func (w *SQLWorker) GetSQLQueryDuration(assignment *Assignment) (time.Duration, error) {
	conn, err := w.ConnectionPool.Acquire(context.Background())
	if err != nil {
		w.Logger.Error("could not aquire connection to database")
		return 0, err
	}
	defer conn.Release()

	// start to count sql query duration
	start := time.Now()
	if err := w.SqlClient.ExecuteQuery(translateAssignmentToQueryBuilder(assignment), conn); err != nil {
		return time.Since(start), errors.Wrap(err, "worker could not execute query")
	}

	pushQueryDurationToPrometheus([]string{"sql"}, float64(time.Since(start)))
	return time.Since(start), nil
}

func pushRequestDurationToPrometheus(labels []string, value float64) {
	requestDurationVec := prompusher.PrometheusClient.NewHistogram("req_duration", "api request duration", labels)
	prompusher.PrometheusClient.PushHistogram(requestDurationVec, "etzba_api_request", "req_duration", labels, value)
}

func pushQueryDurationToPrometheus(labels []string, value float64) {
	queryDurationVec := prompusher.PrometheusClient.NewHistogram("query_duration", "sql query duration", labels)
	prompusher.PrometheusClient.PushHistogram(queryDurationVec, "etzba_sql_query", "query_duration", labels, value)
}
