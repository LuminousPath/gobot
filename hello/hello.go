package hello

import (
	"github.com/mferrera/go-ircevent"
)

func Run(bot *irc.Connection, p, cmd, channel string, word []string, admin bool) {
	say := bot.Privmsg

	// cmd == word[0]
	if cmd == "hello" {
		say(channel, "world!")
	}

	return
}
