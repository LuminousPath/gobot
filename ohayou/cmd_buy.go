package ohayou

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_buy(m common.EmitMsg) {
	if isPM(m) {
		say(m.Nick, "You must buy things in a channel I'm in.")
		return
	}

	if !hasArgs(m.Word) {
		say(m.Channel, "Usage: "+p+"buy <item> will buy you one <item>. "+p+"buy "+
			"<item> 3 will buy you 3 of <item>, if you can afford it.")
		return
	}

	user, ok := GetUser(strings.ToLower(m.Nick))
	if !ok {
		say(m.Channel, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	argOne := strings.ToLower(m.Word[1])
	// just for fun
	if argOne == "ohayou" {
		say(m.Channel, fmt.Sprintf("You purchased %d ohayous for %d ohayous. You "+
			"have %d ohayous left.",
			user.Ohayous, user.Ohayous, user.Ohayous))
		return
	}

	// if a purchase quantity is given
	if len(m.Word) > 2 {
		argTwo := strings.ToLower(m.Word[2])
		// try to convert it to an integer
		amt, err := strconv.Atoi(argTwo)
		if err != nil {
			say(m.Channel, "You didn't give a valid quantity. Usage: "+p+"buy "+
				"<item> will buy you one <item>. "+p+"buy <item> 3 will buy "+
				"you 3 of <item>, if you can afford it.")
			return
		}
		say(m.Channel, user.Buy(m.Channel, argOne, amt))
	} else {
		say(m.Channel, user.Buy(m.Channel, argOne, 1))
	}
}
