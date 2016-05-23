package ohayou

import (
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_register(m common.EmitMsg) {
	var to string
	if isPM(m) {
		to = m.Nick
	} else {
		to = m.Channel
	}
	// if no argument is given for the command
	if !hasArgs(m.Word) {
		say(m.Nick, "Registering allows you to protect your ohayou assets. After "+
			"you are registered, you will be required to identify with the bot "+
			"prior to using most of its commands. Changing your nickname will "+
			"also require you to again identify.")
		say(m.Nick, "Type "+p+"register yes to register your nickname. Your nickname "+
			"must be registered with the network for this to work, and you must "+
			"be identified.")
		say(m.Nick, "To identify with the bot whenever you log on, you must type "+p+
			"identify")
		return
	}

	user, ok := GetUser(strings.ToLower(m.Nick))
	if !ok {
		say(to, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	if strings.ToLower(m.Word[1]) == "yes" {
		if user.Registered {
			say(to, m.Nick+": You are already registered.")
			return
		}
		user.Register(m.Irc)
	}
}
