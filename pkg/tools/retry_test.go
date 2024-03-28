package tools

import (
	"errors"
	"testing"
	"time"
)

func Test_RetryMechanismErr(t *testing.T) {
	r := NewRetry(t, 3, time.Second)
	timestamp := time.Now()

	err := r.Do(func() error {
		return errors.New("test error")
	})

	if err == nil {
		t.Error("expected error")
	}

	if time.Since(timestamp) < 6*time.Second {
		t.Error("expected delay")
	}
}
