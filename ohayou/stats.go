package ohayou

import (
	"fmt"
)

func (u *User) Stats() {
	// if nick is registered but not identified
	if u.Registered && !identified[u.Username] {
		say(u.Username, "You must be identified with me to do that. Make sure "+
			"you are identified with the network and then type "+p+
			"identify.")
		return
	}

	var totalItems int
	var itemsAddOhayous int
	var totalItemsCost int
	var equippedDefense int

	for itm, amt := range u.Items {
		itemMultiplier := 1
		// check if user has item(s) that multiply another item
		if u.ItemMultiply[itm] != 0 {
			itemMultiplier = u.ItemMultiply[itm]
		}
		item, _ := GetItem(itm)
		totalItems += amt
		totalItemsCost += item.Price * amt
		itemsAddOhayous += (item.Add * amt) * itemMultiplier
	}

	for _, itm := range u.Equipped {
		equippedDefense += itm.Defense
	}

	defenseOhayous := equippedDefense / 9
	defenseCats := equippedDefense / 7

	say(u.Username, fmt.Sprintf("You have %d ohayous, having ohayou'd %d times since I "+
		"started keeping track, for which you've received a %d total ohayous. You've "+
		"purchased %d items for %d ohayous and they give you an extra %d ohayous with"+
		" every ohayou.",
		u.Ohayous, u.TimesOhayoued, u.CumOhayous, totalItems, totalItemsCost,
		itemsAddOhayous))

	say(u.Username, fmt.Sprintf("You have %d items equipped and they add %d defense, "+
		"decreasing your chance to have ohayous and cats stolen from by %v%% and "+
		"%v%%, respectively. You've attempted to steal %d times and succeeded %d "+
		"times, having stolen %d ohayous. You've spent %d days on probation for "+
		"being caught stealing, and have had %d ohayous stolen from you.",
		len(u.Equipped), equippedDefense, defenseOhayous, defenseCats,
		u.StealSuccess+u.StealFail, u.StealSuccess, u.StolenOhayous,
		u.ProbationCount, u.OhayousStolen))
}
