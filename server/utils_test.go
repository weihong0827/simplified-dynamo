package main

import (
	"testing"
)

func TestIsKeyInRangeTrue(t *testing.T) {
	var key, start, end uint32
	key = 3
	start = 1
	end = 5

	if !IsKeyInRange(key, start, end) {
		t.Error("Expected true, got false")
	}
}

func TestIsKeyInRangeFalse(t *testing.T) {
	var key, start, end uint32
	key = 3
	start = 5
	end = 1

	if IsKeyInRange(key, start, end) {
		t.Error("Expected false, got true")
	}
}
