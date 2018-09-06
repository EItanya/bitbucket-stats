package logger

import (
	"encoding/json"

	"go.uber.org/zap"
)

var rawJSON = []byte(`{
  "level": "debug",
  "development": true,
  "encoding": "json",
  "outputPaths": ["logs/debug.log"],
  "errorOutputPaths": ["stderr"],
  "encoderConfig": {
    "messageKey": "message",
    "levelKey": "level",
    "levelEncoder": "lowercase"
  }
}`)

// Pointer to main logger
var Log *zap.SugaredLogger

func init() {

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Log = logger.Sugar()

	logger.Info("logger construction succeeded")
}
