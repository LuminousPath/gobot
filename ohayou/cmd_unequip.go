package ohayou

import (
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_unequip(m common.EmitMsg) {
	var to string
	if isPM(m) {
		to = m.Nick
	} else {
		to = m.Channel
	}

	user, ok := GetUser(strings.ToLower(m.Nick))
	// if user has never ohayou'd
	if !ok {
		say(m.Channel, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	// if user didn't provide an item to unequip
	if !hasArgs(m.Word) {
		say(to, "Unequips an equipped item. Usage: "+p+"unequip <item> -- unequips "+
			"<item>")
		return
	}

	// unequip
	say(to, user.Unequip(strings.ToLower(m.Word[1])))

}
