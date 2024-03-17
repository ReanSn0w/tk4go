package tools_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

func Test_RoutineLimiter(t *testing.T) {
	rl := tools.NewRoutineLimiter(1)

	now := time.Now().Unix()
	wg := new(sync.WaitGroup)

	wg.Add(10)

	for i := 0; i < 10; i++ {
		rl.Run(func() {
			if i%2 == 0 {
				time.Sleep(time.Second)
			} else {
				time.Sleep(time.Second * 2)
			}

			t.Log("run")
			wg.Done()
		})
	}

	wg.Wait()
	if time.Now().Unix()-now < 6 {
		t.Error("incorrect work routine limiter")
	}
}
