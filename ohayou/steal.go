package ohayou

import (
	"fmt"
	"math"
)

const (
	// percent chance of success
	stealOhayouSuccess = 36
	stealCatSuccess    = 24

	stealFineMin   int     = 5
	stealFinePct   float64 = 0.16
	stealAmountPct float64 = 0.13
)

func (t *User) stealFine() int {
	return stealFineMin + int(math.Floor(float64(t.Ohayous)*stealFinePct))
}

func (v *User) stealAmount() int {
	return int(math.Floor(float64(v.Ohayous) * stealAmountPct))
}

// t = thief, v = victim
func (t *User) StealFrom(v User, channel, nickRaw, vicRaw string) {
	stealOhayouChance := randNum(0, 100)
	var stealCatChance int
	if v.Items["cat"] == 0 {
		stealCatChance = 0
	} else {
		stealCatChance = randNum(0, 100)
	}

	if t.Ohayous < stealFineMin {
		say(channel, fmt.Sprintf("%s: you don't have enough ohayous to steal. "+
			"You need at least %d.", nickRaw, stealFineMin))
		return
	}

	if v.Ohayous == 0 {
		say(channel, fmt.Sprintf("%s attempts to steal from %s but %s "+
			"doesn't have any ohayous! %s is fined %d ohayous.",
			nickRaw, vicRaw, vicRaw, nickRaw, t.stealFine()))
		return
	}

	// if steal fails on ohayou and cat
	if stealOhayouChance > stealOhayouSuccess && stealCatChance > stealCatSuccess {
		say(channel, fmt.Sprintf("%s attempts to steal from %s but is caught! "+
			"%s is fined %d ohayous",
			nickRaw, vicRaw, nickRaw, t.stealFine()))
		return
	}

	// if steal fails on cat, succeeds on ohayous
	if stealOhayouChance <= stealOhayouSuccess && stealCatChance > stealCatSuccess {
		say(channel, fmt.Sprintf("%s attempts to steal from %s and succeeds! "+
			"%s steals %d ohayous from %s",
			nickRaw, vicRaw, nickRaw, v.stealAmount(), vicRaw))
		return
	}

	// if steal fails on ohayous, succeeds on cat
	if stealOhayouChance > stealOhayouSuccess && stealCatChance <= stealCatSuccess {
		say(channel, fmt.Sprintf("%s attempts to steal from %s and succeeds! "+
			"%s steals a cat from %s",
			nickRaw, vicRaw, nickRaw, vicRaw))
		return
	}

	// if both ohayou and cat thefts succeed
	if stealOhayouChance <= stealOhayouSuccess && stealCatChance <= stealCatSuccess {
		say(channel, fmt.Sprintf("%s attempts to steal from %s and succeeds! "+
			"%s steals a cat and %d ohayous from %s.",
			nickRaw, vicRaw, nickRaw, v.stealAmount(), vicRaw))
		return
	}
}
