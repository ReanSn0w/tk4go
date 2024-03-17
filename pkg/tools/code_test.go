package tools_test

import (
	"testing"
	"time"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

func Test_TimeCodeGeneration(t *testing.T) {
	len := 10
	alive := 10 // код должет будет действовать на временных промежутках до 10 минут
	cg := tools.NewCodeGenerator(alive, len, true, true)

	values := []string{
		"88000001111",
		"79283910412",
		"12308402312",
		"amklasmdio@mail.ru",
		"asldkjowqd@gmail.com",
	}

	for _, v := range values {
		currentTime := time.Time{}
		code := cg.GenerateForTime(currentTime, v)

		for i := 0; i < 60; i++ {
			testTime := currentTime.Add(time.Second * time.Duration(i))
			newCode := cg.GenerateForTime(testTime, v)
			if newCode != code {
				t.Errorf("ct %v - %v, tt %v - %v", currentTime.Unix()%60, currentTime.Unix()/60, testTime.Unix()%60, testTime.Unix()/60)
			}
		}
	}
}

func Test_TimeCodeChecking(t *testing.T) {
	len := 10
	alive := 10 // код должет будет действовать на временных промежутках до 10 минут
	cg := tools.NewCodeGenerator(alive, len, true, true)

	values := []string{
		"88000001111",
		"79283910412",
		"12308402312",
		"amklasmdio@mail.ru",
		"asldkjowqd@gmail.com",
	}

	for _, v := range values {
		currentTime := time.Unix(0, 0)
		code := cg.GenerateForTime(currentTime, v)
		t.Logf("--- Code: %v Time: %v Object: %s", code, currentTime.Unix(), v)
		testTime := currentTime

		for i := 0; i < 30; i++ {
			b := cg.ValidateForTime(testTime, v, code)
			if !b {
				t.Errorf("code testing failed for time %v iteration %v", testTime.Unix(), i)
			}

			testTime = testTime.Add(time.Second * 10)
		}
	}
}
