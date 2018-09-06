package logger

import (
	"go.uber.org/zap"
)

var rawJSON = []byte(`{
  "level": "debug",
  "development": true,
  "encoding": "console",
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

	cfg := zap.NewDevelopmentConfig()
	cfg.Encoding = "console"
	cfg.OutputPaths = []string{"logs/debug.log"}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Log = logger.Sugar()

	logger.Info("logger construction succeeded")
}
