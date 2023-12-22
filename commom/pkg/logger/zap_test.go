package logger

import (
	"google.dev/google/common/pkg/conf"
	"testing"
)

func TestZap(t *testing.T) {
	InitLogger(conf.LoggerConfig{
		Filename: "zaplog.log",
	})

	logger.Info(" sdsadsdadsad")
	logger.Infof(" sdsadsdadsad %d", 123)
	logger.Infof(" sdsadsdadsad %d", 123)
}
