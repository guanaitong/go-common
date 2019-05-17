package alert

import (
	"testing"
	"time"
)

func TestSendByAppName(t *testing.T) {

	SendByAppName(1, "TestSendByAppName")
	time.Sleep(time.Second * 5)

}

func TestSendByCorpCodes(t *testing.T) {
	SendByCorpCodes(7, "TestSendByCorpCodes", "HB266")
	SendByCorpCodes(1, "TestSendByCorpCodes1", "HB266", "HB533")

	time.Sleep(time.Second * 5)
}

func TestSendByGroupId(t *testing.T) {
	SendByGroupId(1, "TestSendByGroupId", 4)
	time.Sleep(time.Second * 5)
}
