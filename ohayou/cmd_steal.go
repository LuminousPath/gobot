package ohayou

import (
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_steal(m common.EmitMsg) {
	if isPM(m) {
		say(m.Nick, "You can only do that in a channel I'm in.")
		return
	}

	// if no argument is given for the command
	if !hasArgs(m.Word) {
		say(m.Channel, "Attempts to steal from someone. Usage: "+p+"steal <user>. Has"+
			" penalties if you are caught!")
		return
	}

	user, ok := GetUser(strings.ToLower(m.Nick))
	// if the thief has never ohayou'd
	if !ok {
		say(m.Channel, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	victim, alsoOk := GetUser(strings.ToLower(m.Word[1]))
	// if the victim has never ohayou'd
	if !alsoOk {
		say(m.Channel, "You can't steal from "+m.Word[1]+" because "+m.Word[1]+
			" has never ohayou'd!")
		return
	}

	go user.StealFrom(victim, m.Channel, m.Nick, m.Word[1], m.Irc)
}
