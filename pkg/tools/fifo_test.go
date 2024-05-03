package tools_test

import (
	"testing"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

func TestFifo_Pop(t *testing.T) {
	f := tools.NewFifo[int]()
	f.Push(1, 2, 3, 4, 5)
	if val := *f.Pop(); val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	if val := *f.Pop(); val != 2 {
		t.Errorf("expected 2, got %d", val)
	}

	if val := *f.Pop(); val != 3 {
		t.Errorf("expected 3, got %d", val)
	}

	if val := *f.Pop(); val != 4 {
		t.Errorf("expected 4, got %d", val)
	}

	if val := *f.Pop(); val != 5 {
		t.Errorf("expected 5, got %d", val)
	}

	if f.Pop() != nil {
		t.Errorf("expected nil, got %d", *f.Pop())
	}
}

func TestFifo_Len(t *testing.T) {
	f := tools.NewFifo[int]()
	f.Push(1, 2, 3, 4, 5)
	if f.Len() != 5 {
		t.Errorf("expected len 5, got %d", f.Len())
	}
}
