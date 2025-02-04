package pokecache

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "https://pokeapi.co/api/v2/location-area/",
			value: []byte("testdata"),
		},
	}

	for _, c := range cases {
		cache := NewCache(5 * time.Second)
		cache.Add(c.key, c.value)

		if _, exists := cache.Get(c.key); !exists {
			t.Errorf("Expected key %s to exist in cache", c.key)
		}

		if val, _ := cache.Get(c.key); string(val) != string(c.value) {
			t.Errorf("Expected value %s, got %s", c.value, val)
		}
	}
}

func TestReadLoop(t *testing.T) {
	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "https://pokeapi.co/api/v2/location-area/",
			value: []byte("testdata"),
		},
	}

	for _, c := range cases {
		interval := 5 * time.Second
		cache := NewCache(interval)

		cache.Add(c.key, c.value)

		time.Sleep(interval + 1*time.Second)

		if _, exists := cache.Get(c.key); exists {
			t.Errorf("Expected not to find ke %s in cache", c.key)
		}
	}
}
