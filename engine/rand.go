package engine

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
)

// ============================================================================
// Random number generator
// This is a seedable RNG for random but repeatable generation of levels etc
// ============================================================================

type gameRNG struct {
	*rand.Rand
}

// Global singleton RNG for the game, for reasons of determinism
var rng gameRNG

func init() {
	s := rand.NewPCG(0, 0)
	rng = gameRNG{rand.New(s)}
}

func seedRNG(seed uint64) {
	s := rand.NewPCG(seed, seed)
	rng = gameRNG{rand.New(s)}
}

func randString(strings ...string) string {
	return strings[rng.IntN(len(strings))]
}

func GetRNG() gameRNG {
	return rng
}

func (r gameRNG) Chance(percentage int) bool {
	return r.IntN(100) < percentage
}

type DiceRoll struct {
	num      int
	sides    int
	modifier int
}

var d100 = DiceRoll{1, 100, 0}

func ParseDiceRoll(dice string) (DiceRoll, bool) {
	re := regexp.MustCompile(`(\d*)d(\d+)([+-]\d+)?`)
	matches := re.FindStringSubmatch(strings.ToLower(dice))

	if len(matches) == 0 {
		return DiceRoll{}, false
	}

	var sides int
	num := 1
	modifier := 0

	if matches[1] != "" {
		num, _ = strconv.Atoi(matches[1])
	}

	if matches[2] != "" {
		sides, _ = strconv.Atoi(matches[2])
	} else {
		return DiceRoll{}, false
	}

	if matches[3] != "" {
		var err error
		modifier, err = strconv.Atoi(matches[3])
		if err != nil {
			return DiceRoll{}, false
		}
	}

	return DiceRoll{num, sides, modifier}, true
}

func (d DiceRoll) Roll() int {
	total := 0
	for i := 0; i < d.num; i++ {
		total += rand.IntN(d.sides) + 1
	}

	return total + d.modifier
}

func (d DiceRoll) String() string {
	if d.modifier == 0 && d.num == 0 && d.sides == 0 {
		return "0"
	}

	if d.num == 1 && d.modifier == 0 {
		return fmt.Sprintf("d%d", d.sides)
	}

	if d.num == 1 {
		return fmt.Sprintf("d%d%+d", d.sides, d.modifier)
	}

	if d.modifier == 0 {
		return fmt.Sprintf("%dd%d", d.num, d.sides)
	}

	return fmt.Sprintf("%dd%d%+d", d.num, d.sides, d.modifier)
}
