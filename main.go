package main

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

type user struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type exampleJSONHandler struct {
	Logger *logrus.Logger
}

func (h exampleJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Start example serve")
	defer h.Logger.Info("End example serve")
	userData := user{Name: "ann", Gender: "female"}
	encoder := json.NewEncoder(w)
	encoder.Encode(userData)
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999999999Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyLevel: "severity",
		},
	})

	// Profiler initialization, best done as early as possible.
	if err := profiler.Start(profiler.Config{
		Service:        "myservice",
		ServiceVersion: "1.0.0",
	}); err != nil {
		// TODO: Handle error.
		logger.Error("Unable to load profiler")
		logger.Error(string(debug.Stack()))
	}

	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		logger.Error(err)
		logger.Error(string(debug.Stack()))
	}
	trace.RegisterExporter(exporter)

	logger.Info("Application Start Up")
	http.Handle("/", exampleJSONHandler{Logger: logger})
	logger.Fatal(http.ListenAndServe("127.0.0.1:8888", &ochttp.Handler{}))
}
