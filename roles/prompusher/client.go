package prompusher

import (
	"os"

	"github.com/nadavbm/zlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"
)

const jobName = "etzba"
const k8sNamespaceFilename = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

type Client interface {
	Set()
	PushGauge(gauge *prometheus.GaugeVec, groupName, groupValue string, labels []string, value float64) error
	PushCounter(counter *prometheus.CounterVec, groupName, groupValue string, labels []string) error
	PushHistogram(histogram *prometheus.HistogramVec, groupName, groupValue string, labels []string, value float64) error
	NewGauge(name, help string, labels []string) *prometheus.GaugeVec
	NewCounter(name, help string, labels []string) *prometheus.CounterVec
	NewHistogram(name, help string, labels []string) *prometheus.HistogramVec
}

// NewClient creates new prometheus api client to push metrics
func NewClient(logger *zlog.Logger) (Client, error) {
	ns, err := getPodNamespace()
	if err != nil {
		return nil, err
	}

	pusher := push.New("asd", jobName)
	return &client{
		Logger:    logger,
		Namespace: ns,
		Pusher:    pusher,
	}, nil
}

type client struct {
	Logger    *zlog.Logger
	Namespace string
	Pusher    *push.Pusher
}

func (c *client) PushGauge(gauge *prometheus.GaugeVec, groupName, groupValue string, labels []string, value float64) error {
	prometheus.MustRegister(gauge)
	gauge.WithLabelValues(labels...).Set(value)
	return c.Pusher.Collector(gauge).Grouping(groupName, groupValue).Push()
}

func (c *client) PushCounter(counter *prometheus.CounterVec, groupName, groupValue string, labels []string) error {
	prometheus.MustRegister(counter)
	counter.WithLabelValues(labels...).Inc()
	return c.Pusher.Collector(counter).Grouping(groupName, groupValue).Push()
}

func (c *client) PushHistogram(histogram *prometheus.HistogramVec, groupName, groupValue string, labels []string, value float64) error {
	prometheus.MustRegister(histogram)
	histogram.WithLabelValues(labels...).Observe(value)
	return c.Pusher.Collector(histogram).Grouping(groupName, groupValue).Push()
}

func getPodNamespace() (string, error) {
	_, inCluster := os.LookupEnv("KUBERNETES_PORT")
	if inCluster {
		bytes, err := os.ReadFile(k8sNamespaceFilename)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
	return "localhost", nil
}

var PrometheusClient Client

func (c *client) Set() {
	c.Logger.Info("setting prometheus client")
}

func Configure(logger *zlog.Logger) error {
	logger.Info("configuring prometheus client")
	var err error
	PrometheusClient, err = NewClient(logger)
	if err != nil {
		logger.Error("could not configure prometheus client", zap.Error(err))
		return err
	}
	PrometheusClient.Set()
	return nil
}
