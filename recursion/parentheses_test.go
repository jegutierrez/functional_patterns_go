package main

import (
	"testing"
)

func TestBalanceParenthesesRecursive(t *testing.T) {

	tt := []struct {
		sentence       string
		expectedResult bool
	}{
		{
			sentence:       "(if (zero? x) max (/ 1 x))",
			expectedResult: true,
		},
		{
			sentence:       "I told him (that it’s not (yet) done). (But he wasn’t listening)",
			expectedResult: true,
		},
		{
			sentence:       ":-)",
			expectedResult: false,
		},
		{
			sentence:       "())(",
			expectedResult: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.sentence, func(t *testing.T) {
			result := BalanceRecursive(tc.sentence)

			if result != tc.expectedResult {
				t.Errorf("unspected result, want: %t, got: %t", tc.expectedResult, result)
			}
		})
	}
}

func TestBalanceParenthesesFor(t *testing.T) {

	tt := []struct {
		sentence       string
		expectedResult bool
	}{
		{
			sentence:       "(if (zero? x) max (/ 1 x))",
			expectedResult: true,
		},
		{
			sentence:       "I told him (that it’s not (yet) done). (But he wasn’t listening)",
			expectedResult: true,
		},
		{
			sentence:       ":-)",
			expectedResult: false,
		},
		{
			sentence:       "())(",
			expectedResult: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.sentence, func(t *testing.T) {
			result := BalanceFor(tc.sentence)

			if result != tc.expectedResult {
				t.Errorf("unspected result, want: %t, got: %t", tc.expectedResult, result)
			}
		})
	}
}

func BenchmarkBalanceParenthesesRecursive(b *testing.B) {
	result := BalanceRecursive(ReallyLargeParenthesesExpression)

	if result != true {
		b.Errorf("unspected result, want true, got: %t", result)
	}
}

func BenchmarkBalanceParenthesesFor(b *testing.B) {
	result := BalanceFor(ReallyLargeParenthesesExpression)

	if result != true {
		b.Errorf("unspected result, want true, got: %t", result)
	}
}
