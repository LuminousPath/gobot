package hello

import (
	"github.com/mferrera/go-ircevent"
)

func Run(b **irc.Connection, p, cmd, channel string, word []string, admin bool) bool {
	bot := *b
	say := bot.Privmsg

	if word[0] == "hello" {
		say(channel, "world!")
	}

	return true
}
