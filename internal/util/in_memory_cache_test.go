package util

import (
	"testing"
	"time"
)

func TestInMemoryCache_Set_Get(t *testing.T) {
	c := NewInMemoryCache[float64](time.Hour)
	c.Set(1.23)
	got, hasValue := c.Get()

	if got != 1.23 {
		t.Fatalf("Get() got = %v, want %v", got, 1.23)
	}
	if !hasValue {
		t.Fatalf("Get() hasValue = %v, want %v", hasValue, true)
	}
}

func TestInMemoryCache_Get_Without_Set(t *testing.T) {
	c := NewInMemoryCache[float64](time.Hour)
	got, hasValue := c.Get()

	if got != 0 {
		t.Fatalf("Get() got = %v, want %v", got, 0)
	}
	if hasValue {
		t.Fatalf("Get() hasValue = %v, want %v", hasValue, false)
	}
}

func TestInMemoryCache_Set_Get_After_TTL(t *testing.T) {
	c := NewInMemoryCache[float64](time.Millisecond)
	c.Set(1.23)
	time.Sleep(2 * time.Millisecond)
	_, hasValue := c.Get()

	if hasValue {
		t.Fatalf("Get() hasValue = %v, want %v", hasValue, false)
	}
}
