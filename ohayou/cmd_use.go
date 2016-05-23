package ohayou

import (
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_use(m common.EmitMsg) {
	if isPM(m) {
		say(m.Nick, "You can only do that in a channel I'm in.")
		return
	}

	// if no argument is given for the command
	if !hasArgs(m.Word) {
		say(m.Channel, "Type "+p+"use <item> to use an item. Type "+p+"inventory to "+
			"see what items you have, or "+p+"items to see what items you can "+p+
			"buy.")
		return
	}

	user, ok := GetUser(strings.ToLower(m.Nick))
	if !ok {
		say(m.Channel, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	argOne := strings.ToLower(m.Word[1])
	if len(m.Word) > 2 {
		argTwo := strings.ToLower(m.Word[2])
		say(m.Channel, user.Use(m.Nick, argOne, argTwo))
	} else {
		say(m.Channel, user.Use(m.Nick, argOne, "somebody"))
	}
}
