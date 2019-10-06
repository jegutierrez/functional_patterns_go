package main

import (
	"testing"
)

func TestCoinsChangeRecursive(t *testing.T) {

	tt := []struct {
		name      string
		amount    int
		coins     []int
		numOfWays int
	}{
		{
			name:      "4 with coins [1,2]",
			amount:    4,
			coins:     []int{1, 2},
			numOfWays: 3,
		},
		{
			name:      "300 with coins [5,10,20,50,100,200,500]",
			amount:    300,
			coins:     []int{5, 10, 20, 50, 100, 200, 500},
			numOfWays: 1022,
		},
		{
			name:      "301 with coins [5,10,20,50,100,200,500] no way",
			amount:    301,
			coins:     []int{5, 10, 20, 50, 100, 200, 500},
			numOfWays: 0,
		},
		{
			name:      "300 with coins [500,5,50,100,20,200,10] unsorted",
			amount:    300,
			coins:     []int{500, 5, 50, 100, 20, 200, 10},
			numOfWays: 1022,
		},
		{
			name:      "1000 with coins [5,10,20,50,100,200,500]",
			amount:    1000,
			coins:     []int{5, 10, 20, 50, 100, 200, 500},
			numOfWays: 104560,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := CoinsChangeRecursive(tc.amount, tc.coins)

			if result != tc.numOfWays {
				t.Errorf("unspected result, want: %d, got: %d", tc.numOfWays, result)
			}
		})
	}
}

func TestCoinsChangeGoStyle(t *testing.T) {

	tt := []struct {
		name      string
		amount    int
		coins     []int
		numOfWays int
	}{
		{
			name:      "4 with coins [1,2]",
			amount:    4,
			coins:     []int{1, 2},
			numOfWays: 3,
		},
		{
			name:      "300 with coins [5,10,20,50,100,200,500]",
			amount:    300,
			coins:     []int{5, 10, 20, 50, 100, 200, 500},
			numOfWays: 1022,
		},
		{
			name:      "301 with coins [5,10,20,50,100,200,500] no way",
			amount:    301,
			coins:     []int{5, 10, 20, 50, 100, 200, 500},
			numOfWays: 0,
		},
		{
			name:      "300 with coins [500,5,50,100,20,200,10] unsorted",
			amount:    300,
			coins:     []int{500, 5, 50, 100, 20, 200, 10},
			numOfWays: 1022,
		},
		{
			name:      "1000 with coins [5,10,20,50,100,200,500]",
			amount:    1000,
			coins:     []int{5, 10, 20, 50, 100, 200, 500},
			numOfWays: 104560,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := CoinsChangeGoStyle(tc.amount, tc.coins)

			if result != tc.numOfWays {
				t.Errorf("unspected result, want: %d, got: %d", tc.numOfWays, result)
			}
		})
	}
}

func BenchmarkCoinsChangeRecursive(b *testing.B) {
	result := CoinsChangeRecursive(3000, []int{5, 10, 20, 50, 100, 200, 500})

	if result != 22481738 {
		b.Errorf("unspected result, want 22481738, got: %d", result)
	}
}

func BenchmarkCoinsChangeGoStyle(b *testing.B) {
	result := CoinsChangeGoStyle(3000, []int{5, 10, 20, 50, 100, 200, 500})

	if result != 22481738 {
		b.Errorf("unspected result, want 22481738, got: %d", result)
	}
}
