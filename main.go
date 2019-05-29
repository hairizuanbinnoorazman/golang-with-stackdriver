package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"cloud.google.com/go/profiler"
	"github.com/sirupsen/logrus"
)

// const LogFile = "/var/log/generic-web-app.log"

// const LogFile = "./generic-web-app.log"

type indexHandler struct{}

func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := "helloworld"
	fmt.Fprint(w, msg)
}

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
	tempValueAlloc := 12
	h.Logger.Warnf("%v", tempValueAlloc)
	h.Logger.Warnf("%+v", userData)
	encoder := json.NewEncoder(w)
	encoder.Encode(userData)
}

func main() {
	// file, _ := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	logger := logrus.New()
	// logger.SetOutput(file)
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
		// ProjectID must be set if not running on GCP.
		// ProjectID: "my-project",
	}); err != nil {
		// TODO: Handle error.
		logger.Error("Unable to load profiler")
		logger.Error(string(debug.Stack()))
	}

	logger.Info("acaccc")
	logger.Error("ckalclacml")
	logger.Warning("ckalclacml")
	// logger.Error(string(debug.Stack()))
	http.Handle("/", exampleJSONHandler{Logger: logger})
	logger.Fatal(http.ListenAndServe("127.0.0.1:8888", nil))
}
