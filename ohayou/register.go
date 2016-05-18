package ohayou

import (
	"fmt"

	"github.com/mferrera/go-ircevent"
)

func (u *User) Register(pin int, b *irc.Connection) {
	var registered bool = false

	b.AddCallback("307", func(e *irc.Event) {
		registered = true
	})
	b.AddCallback("318", func(e *irc.Event) {
		b.ClearCallback("307")
		b.ClearCallback("318")
		if registered {
			say(u.Username, fmt.Sprintf("Successfully registered your PIN: %d. "+
				"Don't forget it!", pin))
			u.SavePin(pin)
		} else {
			say(u.Username, "Your nick isn't registered or you aren't identified "+
				"with NickServ. You must do both before you can register")
		}
		return

	})
	b.Whois(u.Username)
}
