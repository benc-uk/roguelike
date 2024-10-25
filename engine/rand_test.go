package engine

// ============================================================================
// Tests for dice rolls
// ============================================================================

import (
	"testing"
)

func TestParseDiceRoll(t *testing.T) {
	tests := []struct {
		input    string
		expected DiceRoll
		valid    bool
	}{
		{"d6", DiceRoll{1, 6, 0}, true},
		{"2d8+3", DiceRoll{2, 8, 3}, true},
		{"3d10-2", DiceRoll{3, 10, -2}, true},
		{"d12", DiceRoll{1, 12, 0}, true},
		{"4d", DiceRoll{}, false},
		{"2d6+abc", DiceRoll{2, 6, 0}, true},
		{"17D20+150", DiceRoll{17, 20, 150}, true},
		{"", DiceRoll{}, false},
		{"dice", DiceRoll{}, false},
	}

	for _, test := range tests {
		result, valid := ParseDiceRoll(test.input)
		if valid != test.valid || result != test.expected {
			t.Errorf("ParseDiceRoll(%q) = %v, %v; want %v, %v", test.input, result, valid, test.expected, test.valid)
		}
	}
}

func TestDiceRoll_Roll(t *testing.T) {
	tests := []struct {
		dice     DiceRoll
		min, max int
	}{
		{DiceRoll{1, 6, 0}, 1, 6},
		{DiceRoll{2, 8, 3}, 5, 19},
		{DiceRoll{3, 10, -2}, 1, 28},
		{DiceRoll{1, 12, 0}, 1, 12},
		{DiceRoll{17, 20, 150}, 167, 490},
		{DiceRoll{0, 0, 0}, 0, 0},
	}

	for _, test := range tests {
		for i := 0; i < 100; i++ {
			result := test.dice.Roll()
			if result < test.min || result > test.max {
				t.Errorf("DiceRoll(%v).Roll() = %d; want between %d and %d", test.dice, result, test.min, test.max)
			}
		}
	}
}
