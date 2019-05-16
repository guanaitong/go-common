package logs

import (
	"github.com/golang/glog"
	"testing"
)

func TestInfo(t *testing.T) {
	glog.Info("OK")
}

func TestInfo2(t *testing.T) {
	InitLogs()
	defer FlushLogs()

	glog.Info("OK")

	logger := NewLogger("GAT ")
	logger.Println("logger")
}
