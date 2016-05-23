package ohayou

import (
	"fmt"
	"math"
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_ohayous(m common.EmitMsg) {
	user, ok := GetUser(strings.ToLower(m.Nick))
	if !ok {
		say(m.Channel, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	// if nick is registered but not identified
	if user.Registered && !identified[user.Username] {
		say(user.Username, user.Username+": You must be identified with me to do "+
			"that. Make sure "+"you are identified with the network and then "+
			"type "+p+"identify.")
		return
	}

	if user.Vault.Installed {
		say(m.Nick, fmt.Sprintf("You have %d ohayous on hand and %d ohayous in your "+
			"Level %d vault. Your vault's capacity is %d ohayous.",
			user.Ohayous, user.Vault.Ohayous, user.Vault.Level+1,
			int(math.Pow(10, 2+float64(user.Vault.Level)))))
	} else {
		say(m.Nick, fmt.Sprintf("You have %d ohayous.", user.Ohayous))
	}
}
