package zap

import (
	"bytes"
	"context"
	"sync"

	"github.com/go-kita/log/v3"

	"go.uber.org/zap"
)

// output is a log.Output based on Zap logger.
type output struct {
	out *zap.Logger
	buf *sync.Pool
}

// NewOutput produce a log.Output based on zap.Logger
func NewOutput(out *zap.Logger) log.Output {
	return &output{
		out: out,
		buf: &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

func (o *output) Output(ctx context.Context, level log.Level, msg string, fields []log.Field, addCallerSkip int) {
	o.outFunc(level, addCallerSkip)(msg, zapFields(ctx, fields)...)
}

type outFunc func(msg string, fields ...zap.Field)

func (o *output) outFunc(level log.Level, addCallerSkip int) outFunc {
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
	case level == log.ClosedLevel:
		return func(_ string, _ ...zap.Field) {
		}
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
	return log.NewStdLogger(name, NewOutput(out))
}

// MakeProvider make a log.LoggerProvider function.
func MakeProvider(out *zap.Logger) log.LoggerProvider {
	return func(name string) log.Logger {
		return NewLogger(name, out)
	}
}
