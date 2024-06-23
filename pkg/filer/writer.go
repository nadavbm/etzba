package filer

import (
	"encoding/json"
	"os"

	"github.com/nadavbm/zlog"
	"go.uber.org/zap"
)

// Writer used for writing results files
type Writer struct {
	logger *zlog.Logger
}

// Writer creates an instance of a writer
func NewWriter(logger *zlog.Logger) *Writer {
	return &Writer{
		logger: logger,
	}
}

func (w *Writer) WriteFile(file string, obj interface{}) error {
	content, err := json.MarshalIndent(&obj, "", "    ")
	if err != nil {
		return err
	}

	w.logger.Info("writing content to", zap.String("file", file))
	return os.WriteFile(file, content, 0644)
}
