package lib

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logger *zap.SugaredLogger
}

func GormLoggerNew(logger *zap.SugaredLogger) *GormLogger {
	return &GormLogger{logger: logger}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	g.logger.Infof(s, i)
}

func (g *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	g.logger.Warnf(s, i)
}

func (g *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	g.logger.Errorf(s, i)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

}
