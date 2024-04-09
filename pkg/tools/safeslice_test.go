package tools_test

import (
	"testing"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

func TestSafeSlice(t *testing.T) {
	// Create a new SafeSlice
	ss := tools.NewSafeSlice[int]()
	// Push elements into the SafeSlice
	ss.Push(1, 2, 3)

	// Test the Len method
	if ss.Len() != 3 {
		t.Errorf("Expected length 3, got %d", ss.Len())
	}

	// Test the Get method
	if ss.Get(0) != 1 {
		t.Errorf("Expected element at index 0 to be 1, got %d", ss.Get(0))
	}

	// Test the Delete method
	ss.Delete(1)
	if ss.Len() != 2 {
		t.Errorf("Expected length 2 after deleting an element, got %d", ss.Len())
	}
	if ss.Get(1) != 3 {
		t.Errorf("Expected element at index 1 to be 3 after deleting an element, got %d", ss.Get(1))
	}
}
