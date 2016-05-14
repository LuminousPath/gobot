package ohayou

import (
	"fmt"

	"github.com/mferrera/go-ircevent"
)

var (
	isRegistered bool
)

func doRegister(nick string, pn int) {
	b.AddCallback("307", func(e *irc.Event) {
		isRegistered = true
	})
	b.AddCallback("318", func(e *irc.Event) {
		b.ClearCallback("307")
		b.ClearCallback("318")
		if isRegistered {
			if getUser(nick) {
				say(nick, fmt.Sprintf("Successfully registered your PIN: %d. "+
					"Don't forget it!", pn))
				USER.savePin(pn)
			} else {
				say(nick, "Looks like you haven't ohayou'd yet. Type "+p+
					"ohayou in a channel I'm in to get your ration, and "+
					"then you can register.")
			}
		} else {
			say(nick, "Your nick isn't registered or you aren't identified with "+
				"NickServ. You must do both before you can register")
		}
		isRegistered = false
		return

	})
	b.Whois(nick)
}
