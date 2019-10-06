package main

import (
	"testing"
)

func TestBalanceParenthesesTailRecursion(t *testing.T) {

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
			result := BalanceTailRecursive(tc.sentence)

			if result != tc.expectedResult {
				t.Errorf("unspected result, want: %t, got: %t", tc.expectedResult, result)
			}
		})
	}
}

func TestBalanceParenthesesGoStyle(t *testing.T) {

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
			result := BalanceGoStyle(tc.sentence)

			if result != tc.expectedResult {
				t.Errorf("unspected result, want: %t, got: %t", tc.expectedResult, result)
			}
		})
	}
}

func BenchmarkBalanceParenthesesTailRecursion(b *testing.B) {
	result := BalanceTailRecursive(ReallyLargeParenthesesExpression)

	if result != true {
		b.Errorf("unspected result, want true, got: %t", result)
	}
}

func BenchmarkBalanceParenthesesGoStyle(b *testing.B) {
	result := BalanceGoStyle(ReallyLargeParenthesesExpression)

	if result != true {
		b.Errorf("unspected result, want true, got: %t", result)
	}
}
