package gamification

import "testing"

func TestXPForLevel(t *testing.T) {

	tests := []struct {
		level int
		want  int
	}{
		{1, 0},
		{2, 100},
		{3, 300},
		{4, 600},
		{5, 1000},
	}

	for _, test := range tests {
		if got := XPForLevel(test.level); got != test.want {
			t.Errorf("XPForLevel(%d) = %d, want %d", test.level, got, test.want)
		}
	}
}

func TestLevelForXP(t *testing.T) {

	tests := []struct {
		xp                 int
		wantLevel          int
		wantCurrentLevelXP int
		wantXPForNextLevel int
	}{
		{0, 1, 0, 100},
		{50, 1, 50, 100},
		{99, 1, 99, 100},
		{100, 2, 0, 200},
		{299, 2, 199, 200},
		{300, 3, 0, 300},
		{999, 4, 399, 400},
	}

	for _, test := range tests {

		level, currentLevelXP, xpForNextLevel := LevelForXP(test.xp)

		if level != test.wantLevel {
			t.Errorf("LevelForXP(%d) level = %d, want %d", test.xp, level, test.wantLevel)
		}

		if currentLevelXP != test.wantCurrentLevelXP {
			t.Errorf("LevelForXP(%d) currentLevelXP = %d, want %d", test.xp, currentLevelXP, test.wantCurrentLevelXP)
		}

		if xpForNextLevel != test.wantXPForNextLevel {
			t.Errorf("LevelForXP(%d) xpForNextLevel = %d, want %d", test.xp, xpForNextLevel, test.wantXPForNextLevel)
		}
	}
}
