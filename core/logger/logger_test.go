package logger

import (
	"context"
	"github.com/tonly18/xws/core/zutils"
	"testing"
)

func init() {
	Init()
}

func TestLogger(t *testing.T) {
	ctx := context.WithValue(context.Background(), zutils.TraceID, "trace_id_1234567890")
	Infof(ctx, "TestLogger:%s", "testlogger-error")
	abc(ctx)
}

func abc(ctx context.Context) {
	Infof(ctx, "abc error:%s", "abc-error")
	def(ctx)
}

func def(ctx context.Context) {
	Infof(ctx, "def error:%s", "def-error")
	ghi(ctx)
}

func ghi(ctx context.Context) {
	Infof(ctx, "ghi error:%s", "ghi-error")
}
