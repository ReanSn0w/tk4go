package tools_test

import (
	"testing"
	"time"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

func TestScheduler_Wait(t *testing.T) {
	tasks := []int{1, 2, 3}
	done := map[int]bool{
		1: false,
		2: false,
		3: false,
	}

	s := tools.NewScheduler[int, int](2, func(i int) int {
		time.Sleep(time.Second * time.Duration(5-i))
		return i
	})

	s.Push(tasks...)

	go func() {
		for i := range s.Out() {
			t.Logf("got %d", i)
			done[i] = true
		}
	}()

	s.Wait()
}
