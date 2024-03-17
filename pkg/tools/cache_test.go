package tools_test

import (
	"testing"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

var (
	cacheValues = []struct {
		Key   string
		Value int
	}{
		{Key: "first", Value: 10},
		{Key: "second", Value: 1234},
		{Key: "thrid", Value: 1283},
	}
)

func Test_CacheCRUD(t *testing.T) {
	cache := tools.NewCache[int]()

	// Set
	for _, pair := range cacheValues {
		cache.Set(pair.Key, &pair.Value)
	}

	// Get
	for _, pair := range cacheValues {
		if &pair.Value != cache.Get(pair.Key) {
			t.Log("unvalid pair for key: ", pair.Key)
		}
	}

	// Delete
	for _, pair := range cacheValues {
		cache.Delete(pair.Key)
	}

	for _, pair := range cacheValues {
		if cache.Get(pair.Key) != nil {
			t.Log("pair value found after delete. key: ", pair.Key)
		}
	}
}
