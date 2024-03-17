package tools_test

import (
	"context"
	"testing"
	"time"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

func Test_Shutdown(t *testing.T) {
	st := tools.NewShutdownStack(tools.BaseLogger())
	ctx, cancel := context.WithCancel(context.Background())

	count := 0
	st.Add(func(ctx context.Context) {
		count++
	}, func(ctx context.Context) {
		count++
	}, func(ctx context.Context) {
		count++
	})

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	st.Wait(ctx, time.Second*10)
	if count != 3 {
		t.Error("count not equal 3")
	}
}
