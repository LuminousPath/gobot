package hello

import (
	"github.com/mferrera/go-ircevent"
)

func Run(bot *irc.Connection, p, cmd, channel string, word []string, admin bool) {
	say := bot.Privmsg

	if word[0] == "hello" {
		say(channel, "world!")
	}

	return
}
