package zap

import (
	"context"

	"github.com/go-kita/log"
	"go.uber.org/zap"
)

// outPutter is a log.OutPutter based on Zap logger.
type outPutter struct {
	out *zap.Logger
}

var _ log.OutPutter = (*outPutter)(nil)

// NewOutPutter produce a log.OutPutter based on zap.Logger
func NewOutPutter(out *zap.Logger) log.OutPutter {
	return &outPutter{
		out: out,
	}
}

func (o *outPutter) OutPut(
	ctx context.Context, _ string, level log.Level, msg string, fields []log.Field, callDepth int) {
	o.outFunc(level, callDepth)(msg, zapFields(ctx, fields)...)
}

type outFunc func(msg string, fields ...zap.Field)

func (o *outPutter) outFunc(level log.Level, addCallerSkip int) outFunc {
	out := o.out
	out = out.WithOptions(zap.AddCallerSkip(addCallerSkip + 2))
	switch {
	case level <= log.DebugLevel:
		return out.Debug
	case level == log.InfoLevel:
		return out.Info
	case level == log.WarnLevel:
		return out.Warn
	case level == log.ErrorLevel:
		return out.Error
	default:
		return out.Error
	}
}

func zapFields(ctx context.Context, fields []log.Field) []zap.Field {
	n := len(fields)
	if n == 0 {
		return nil
	}
	zfs := make([]zap.Field, 0, n)
	for _, f := range fields {
		if f.Key == log.LevelKey {
			continue
		}
		zfs = append(zfs, zap.Any(f.Key, log.Value(ctx, f.Value)))
	}
	return zfs
}

// NewLogger create a log.Logger of name based on zap.Logger.
func NewLogger(name string, out *zap.Logger) log.Logger {
	return log.NewStdLogger(name, NewOutPutter(out))
}
