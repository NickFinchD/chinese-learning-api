package review

import (
	"testing"
	"time"
)

func TestCalculateNextReview(t *testing.T) {

	tests := []struct {
		name         string
		correctCount int
		correct      bool
		expectedDays int
	}{
		{
			name:         "first correct answer",
			correctCount: 0,
			correct:      true,
			expectedDays: 1,
		},
		{
			name:         "second correct answer",
			correctCount: 1,
			correct:      true,
			expectedDays: 3,
		},
		{
			name:         "third correct answer",
			correctCount: 2,
			correct:      true,
			expectedDays: 7,
		},
		{
			name:         "fourth correct answer",
			correctCount: 3,
			correct:      true,
			expectedDays: 14,
		},
		{
			name:         "fifth correct answer",
			correctCount: 4,
			correct:      true,
			expectedDays: 30,
		},
		{
			name:         "wrong answer",
			correctCount: 10,
			correct:      false,
			expectedDays: 0,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			progress := &UserWordProgress{
				CorrectCount: test.correctCount,
			}

			next := calculateNextReview(progress, test.correct)

			expected := time.Now().Add(time.Duration(test.expectedDays) * 24 * time.Hour)

			diff := next.Sub(expected)

			if diff > time.Second || diff < -time.Second {
				t.Fatalf(
					"expected about %d days, got diff %v",
					test.expectedDays,
					diff,
				)
			}
		})
	}
}
