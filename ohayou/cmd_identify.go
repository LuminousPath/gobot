package ohayou

import (
	"strings"

	"github.com/mferrera/go-ircevent"
	"github.com/mferrera/gobot/common"
)

func cmd_identify(m common.EmitMsg) {
	var to string
	if isPM(m) {
		to = m.Nick
	} else {
		to = m.Channel
	}
	user, ok := GetUser(strings.ToLower(m.Nick))
	if !ok {
		say(to, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	if !user.Registered {
		say(to, "You can't do that because you're not registered yet! Type "+p+
			"register to be PM'd information about registering.")
		return
	}

	if identified[user.Username] {
		say(to, user.Username+": You are already identified.")
		return
	}

	var networkReg bool = false

	m.Irc.AddCallback("307", func(e *irc.Event) {
		networkReg = true
	})
	m.Irc.AddCallback("318", func(e *irc.Event) {
		m.Irc.ClearCallback("307")
		m.Irc.ClearCallback("318")
		if networkReg {
			say(to, user.Username+": You are now identified with the bot. "+
				"Changing nicks or logging off will remove this.")
			if !watchingNicks {
				watchNicks(m.Irc)
			}
			identified[user.Username] = true
		} else {
			say(to, user.Username+": You must be identified with the network to "+
				"identify with me.")
		}
		return

	})
	m.Irc.Whois(user.Username)
}

func watchNicks(b *irc.Connection) {
	b.AddCallback("NICK", func(e *irc.Event) {
		if identified[strings.ToLower(e.Nick)] {
			delete(identified, strings.ToLower(e.Nick))
		}
	})
	b.AddCallback("QUIT", func(e *irc.Event) {
		if identified[strings.ToLower(e.Nick)] {
			delete(identified, strings.ToLower(e.Nick))
		}
	})
}
