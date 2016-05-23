package ohayou

import (
	"github.com/mferrera/go-ircevent"
)

func (u *User) Register(b *irc.Connection) {
	var networkReg bool = false

	b.AddCallback("307", func(e *irc.Event) {
		networkReg = true
	})
	b.AddCallback("318", func(e *irc.Event) {
		b.ClearCallback("307")
		b.ClearCallback("318")
		if networkReg {
			say(u.Username, "Successfully registered! Type "+p+"identify to "+
				"identify yourself with the bot.")
			u.SetRegister(true)
		} else {
			say(u.Username, "Your nick isn't registered or you aren't identified "+
				"with NickServ. You must do both before you can register")
		}
		return

	})
	b.Whois(u.Username)
}
