package task

import (
	"fmt"
	"testing"
	"time"
)

func TestStartBackgroundTask(t *testing.T) {

	count := 0
	StartBackgroundTask("test1", time.Second, func() {
		count = count + 1
		fmt.Println("123")
	})

	time.Sleep(time.Second * 5)

	if count < 3 {
		t.Errorf("not work")
	}
}

func TestStartBackgroundTaskCrash(t *testing.T) {
	count := 0
	StartBackgroundTask("test1", time.Second, func() {
		count = count + 1
		panic("123")
	})
	time.Sleep(time.Second * 5)
	if count < 3 {
		t.Errorf("not work")
	}
}
