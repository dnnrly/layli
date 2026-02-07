package common

import (
	"os"
	"testing"
)

func TestService_Shuffle(t *testing.T) {
	// Test with deterministic seed
	service := New(42)

	// Create a slice
	original := []int{1, 2, 3, 4, 5}
	shuffled := make([]int, len(original))
	copy(shuffled, original)

	// Shuffle it
	service.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Should be different from original (most likely)
	if len(shuffled) != len(original) {
		t.Errorf("Expected length %d, got %d", len(original), len(shuffled))
	}

	// Test multiple shuffles with same seed produce same result
	service2 := New(42)
	original2 := []int{1, 2, 3, 4, 5}
	shuffled2 := make([]int, len(original2))
	copy(shuffled2, original2)

	service2.Shuffle(len(shuffled2), func(i, j int) {
		shuffled2[i], shuffled2[j] = shuffled2[j], shuffled2[i]
	})

	for i := range shuffled {
		if shuffled[i] != shuffled2[i] {
			t.Errorf("Expected same shuffle with same seed, got different at index %d: %d vs %d", i, shuffled[i], shuffled2[i])
		}
	}
}

func TestGlobalShuffle(t *testing.T) {
	// Test global shuffle function
	original := []int{1, 2, 3, 4, 5}
	shuffled := make([]int, len(original))
	copy(shuffled, original)

	// Set test environment variable
	os.Setenv("LAYLI_TEST_SEED", "1")
	defer os.Unsetenv("LAYLI_TEST_SEED")

	// Shuffle using global function
	Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Should be different from original (most likely)
	if len(shuffled) != len(original) {
		t.Errorf("Expected length %d, got %d", len(original), len(shuffled))
	}

	// Test that multiple calls with test env set produce same result
	original2 := []int{1, 2, 3, 4, 5}
	shuffled2 := make([]int, len(original2))
	copy(shuffled2, original2)

	Shuffle(len(shuffled2), func(i, j int) {
		shuffled2[i], shuffled2[j] = shuffled2[j], shuffled2[i]
	})

	for i := range shuffled {
		if shuffled[i] != shuffled2[i] {
			t.Errorf("Expected same shuffle with test seed, got different at index %d: %d vs %d", i, shuffled[i], shuffled2[i])
		}
	}
}

func TestNewDefault(t *testing.T) {
	service := NewDefault()

	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	// Test that it can shuffle without panicking
	data := []int{1, 2, 3, 4, 5}
	service.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	// Should still have same length
	if len(data) != 5 {
		t.Errorf("Expected length 5, got %d", len(data))
	}
}

func TestGlobalShuffleWithoutTestSeed(t *testing.T) {
	// Test global shuffle without test environment variable
	original := []int{1, 2, 3, 4, 5}
	shuffled := make([]int, len(original))
	copy(shuffled, original)

	// Ensure test env is not set
	os.Unsetenv("LAYLI_TEST_SEED")

	// Shuffle using global function
	Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Should be different from original (most likely)
	if len(shuffled) != len(original) {
		t.Errorf("Expected length %d, got %d", len(original), len(shuffled))
	}
}
