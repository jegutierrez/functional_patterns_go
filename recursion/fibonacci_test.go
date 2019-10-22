package main

import (
	"testing"
)

func TestFibonacciRecursive(t *testing.T) {

	tt := []struct {
		name   string
		n      int
		result int
	}{
		{
			name:   "fib of 4",
			n:      4,
			result: 3,
		},
		{
			name:   "fib of 10",
			n:      10,
			result: 55,
		},
		{
			name:   "fib of 100",
			n:      45,
			result: 1134903170,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := FibonacciRecursive(tc.n)

			if result != tc.result {
				t.Errorf("unspected result, want: %d, got: %d", tc.result, result)
			}
		})
	}
}

func TestFibonacciFor(t *testing.T) {

	tt := []struct {
		name   string
		n      int
		result int
	}{
		{
			name:   "fib of 4",
			n:      4,
			result: 3,
		},
		{
			name:   "fib of 10",
			n:      10,
			result: 55,
		},
		{
			name:   "fib of 100",
			n:      45,
			result: 1134903170,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := FibonacciFor(tc.n)

			if result != tc.result {
				t.Errorf("unspected result, want: %d, got: %d", tc.result, result)
			}
		})
	}
}

func BenchmarkFibonacciRecursive(b *testing.B) {
	result := FibonacciRecursive(45)

	if result != 1134903170 {
		b.Errorf("unspected result, want 22481738, got: %d", result)
	}
}

func BenchmarkFibonacciFor(b *testing.B) {
	result := FibonacciFor(45)

	if result != 1134903170 {
		b.Errorf("unspected result, want 22481738, got: %d", result)
	}
}
