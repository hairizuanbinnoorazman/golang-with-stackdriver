package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"cloud.google.com/go/profiler"
	"github.com/sirupsen/logrus"
)

const ServiceName = "external-service"
const ServiceHost = "0.0.0.0"
const ServicePort = 8888
const ServiceVersion = "v0.0.1"

type errResponse struct {
	Error string `json:"error"`
}
type indexHandler struct {
	Logger *logrus.Logger
}

func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Start Index Handler")
	defer h.Logger.Info("End Index handler")
	type respIndexHandler struct {
		Status string `json:"status"`
	}
	resp := respIndexHandler{Status: "Ok"}
	respOutput, _ := json.Marshal(resp)
	w.Write(respOutput)
}

type logTestHandler struct {
	Logger   *logrus.Logger
	Severity string
}

func (h logTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Start Log Test Handler")
	defer h.Logger.Info("End Log Test handler")
	switch h.Severity {
	case "info":
		h.Logger.Info("Info log from Log Test Handler")
	case "warning":
		h.Logger.Warning("Warning log from Log Test Handler")
	case "error":
		h.Logger.Error("Error log from Log Test Handler")
	case "stack":
		h.Logger.Error("Error log from Log Test Handler")
		h.Logger.Error(string(debug.Stack()))
	}
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
		Service:        ServiceName,
		ServiceVersion: ServiceVersion,
	}); err != nil {
		// TODO: Handle error.
		logger.Errorf("Unable to load profiler %v", err)
		logger.Error(string(debug.Stack()))
	}

	logger.Info("Application Start Up")
	http.Handle("/", indexHandler{Logger: logger})
	http.Handle("/info", logTestHandler{Logger: logger, Severity: "info"})
	http.Handle("/warning", logTestHandler{Logger: logger, Severity: "warning"})
	http.Handle("/error", logTestHandler{Logger: logger, Severity: "error"})
	http.Handle("/stack", logTestHandler{Logger: logger, Severity: "stack"})
	logger.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", ServiceHost, ServicePort), nil))
}
