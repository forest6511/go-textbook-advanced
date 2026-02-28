package main_test

import (
	"testing"
	"testing/synctest"
	"time"
)

// Table-Driven Test example
func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"positive", 1, 2, 3},
		{"zero", 0, 0, 0},
		{"negative", -1, -2, -3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a + tc.b
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

// Cache for synctest example
type Cache struct {
	value   string
	expires time.Time
}

func (c *Cache) Get() (string, bool) {
	if time.Now().After(c.expires) {
		return "", false
	}
	return c.value, true
}

// TestCacheExpiry uses testing/synctest (Go 1.25)
func TestCacheExpiry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		now := time.Now()
		c := &Cache{value: "hello", expires: now.Add(5 * time.Second)}

		if _, ok := c.Get(); !ok {
			t.Error("expected cache hit before expiry")
		}

		time.Sleep(10 * time.Second)

		if _, ok := c.Get(); ok {
			t.Error("expected cache miss after expiry")
		}
	})
}

// BenchmarkLoop uses b.Loop() (Go 1.24+)
func BenchmarkLoop(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	for b.Loop() {
		sum := 0
		for _, v := range data {
			sum += v
		}
		_ = sum
	}
}

// TestArtifactDir demonstrates t.ArtifactDir (Go 1.26)
func TestArtifactDir(t *testing.T) {
	dir := t.ArtifactDir()
	if dir == "" {
		t.Error("ArtifactDir returned empty string")
	}
	t.Logf("artifact dir: %s", dir)
}
