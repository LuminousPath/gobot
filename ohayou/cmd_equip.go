package ohayou

import (
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_equip(m common.EmitMsg) {
	var to string
	if isPM(m) {
		to = m.Nick
	} else {
		to = m.Channel
	}

	user, ok := GetUser(strings.ToLower(m.Nick))
	// user has never ohayou'd
	if !ok {
		say(to, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	// user didnt provide an item to equip
	if !hasArgs(m.Word) {
		say(to, "Type "+p+"equip <item> to equip <item>. You can only have one item "+
			"equipped per slot, unless otherwise noted.")
		return
	}

	// equip item and say result
	say(to, user.Equip(strings.ToLower(m.Word[1])))
}
