package tools_test

import (
	"testing"
	"time"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

func Test_CycleOnce(t *testing.T) {
	value := 0

	c := tools.NewCycleTask(func() {
		value += 1
	})

	for i := 1; i <= 5; i++ {
		c.Once()

		if value != i {
			t.Log("once execution failed")
			t.Fail()
		}
	}
}

func Test_CycleSecond(t *testing.T) {
	value := 0

	c := tools.NewCycleTask(func() {
		value += 1
	})

	c.Run(time.Second)
	time.Sleep(time.Second * 5)

	if value < 4 {
		t.Log("timer operation failed")
		t.Fail()
	}

	c.Stop()
	time.Sleep(time.Second * 3)

	if value > 7 {
		t.Log("stop execution failed")
		t.Fail()
	}
}

func Test_CycleStopTimer(t *testing.T) {
	value := 0

	c := tools.NewCycleTask(func() {
		value += 1
	})

	c.Run(time.Second * 10)
	time.Sleep(time.Second * 2)
	c.Stop()

	if value != 0 {
		t.Log("stoppping timer failed")
		t.Fail()
	}

}
