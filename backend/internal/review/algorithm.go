package review

import "time"

func calculateNextReview(
	progress *UserWordProgress,
	correct bool,
) time.Time {

	if !correct {
		return time.Now()
	}

	switch progress.CorrectCount {

	case 0:
		return time.Now().Add(24 * time.Hour)

	case 1:
		return time.Now().Add(3 * 24 * time.Hour)

	case 2:
		return time.Now().Add(7 * 24 * time.Hour)

	case 3:
		return time.Now().Add(14 * 24 * time.Hour)

	default:
		return time.Now().Add(30 * 24 * time.Hour)
	}
}
