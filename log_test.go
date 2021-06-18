package zap

import (
	"bytes"
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/go-kita/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MemorySink struct {
	*bytes.Buffer
}

func (m *MemorySink) Sync() error {
	return nil
}

func (m *MemorySink) Close() error {
	return nil
}

func TestZapLog(t *testing.T) {
	buf := &bytes.Buffer{}
	sink := &MemorySink{buf}
	_ = zap.RegisterSink("memory", func(url *url.URL) (zap.Sink, error) {
		return sink, nil
	})

	conf := zap.NewDevelopmentConfig()
	conf.OutputPaths = []string{"memory://"}
	conf.Level.SetLevel(zapcore.DebugLevel)
	zl, err := conf.Build(zap.AddCaller())
	if err != nil {
		t.Fatalf("error creating zap logger")
	}
	log.GetLevelStore().Set("", log.DebugLevel)
	log.GetLevelStore().Set("closed", log.ClosedLevel)
	logger := NewLogger("", zl)
	for i := log.DebugLevel; i <= log.ErrorLevel; i++ {
		logger.AtLevel(context.Background(), i).
			Print("abc")
		got := buf.String()
		if !strings.Contains(got, "_test.go:") {
			t.Errorf("expect %s outPutter contains current test file name, but got: %q",
				i, got)
		}
		buf.Reset()
	}
	closed := NewLogger("closed", zl)
	closed.AtLevel(context.Background(), log.InfoLevel).Print("abc")
	got := buf.String()
	buf.Reset()
	if got != "" {
		t.Errorf("expect closed outPutter nothing, but got %q", got)
	}
}
