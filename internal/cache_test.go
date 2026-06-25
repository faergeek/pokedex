package internal

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	entries := []struct{ key, value string }{
		{key: "https//example.com/users/1", value: "foo"},
		{key: "https//example.com/users/3", value: "bar"},
	}

	cache := NewCache(time.Millisecond)
	for i := range entries {
		cache.Add(entries[i].key, []byte(entries[i].value))
	}

	for i := range entries {
		expected := entries[i].value

		actual, exists := cache.Get(entries[i].key)
		if !exists {
			t.Errorf("expected cache entry to exist right after it's creation")
		}

		if len(actual) != len(expected) {
			t.Errorf("expected length to be %d, got %d", len(actual), len(expected))
		}

		for i := range actual {
			if actual[i] != expected[i] {
				t.Errorf("expected byte %d to be %d, got %d", i, expected[i], actual[i])
			}
		}
	}
}

func TestGetNoEntry(t *testing.T) {
	cache := NewCache(time.Millisecond)
	_, exists := cache.Get("foo")
	if exists {
		t.Errorf("expected cache entry to not exist since it was not added")
	}
}

func TestReapLoop(t *testing.T) {
	lifetime := time.Millisecond

	cache := NewCache(lifetime)
	cache.Add("foo", []byte("hello"))

	_, exists := cache.Get("foo")
	if !exists {
		t.Errorf("expected cache entry to exist right after it's creation")
	}

	time.Sleep(lifetime + 2*time.Millisecond)

	_, exists = cache.Get("foo")
	if exists {
		t.Errorf("expected cache entry to not exist anymore after it's lifetime has ended")
	}
}
