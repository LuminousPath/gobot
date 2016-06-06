package ohayou

import (
	"github.com/mferrera/gobot/common"
)

func cmd_help(m common.EmitMsg) {
	var to string
	if isPM(m) {
		to = m.Nick
	} else {
		to = m.Channel
	}

	say(to, "An ohayou game. Acquire ohayous and purchase things with them. Some items "+
		"have special functions. Commands: "+p+"ohayou, "+p+"buy, "+p+"item, "+p+
		"items, "+p+"use, "+p+"inventory, "+p+"stats, "+p+"equip, "+p+
		"register, "+p+"changelog.")
}
