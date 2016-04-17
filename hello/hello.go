package hello

import (
	"github.com/thoj/go-ircevent"
)

func Run(b **irc.Connection, p, cmd, channel string, word []string, admin bool) {
	bot := *b
	say := bot.Privmsg

	if word[0] == "hello" {
		say(channel, "world!")
	}

	return
}
