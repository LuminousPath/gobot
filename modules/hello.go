package hello

import (
	"github.com/thoj/go-ircevent"
)

func Run(b **irc.Connection, p, cmd, channel string, word []string, admin bool) {
	bot := *b

	if word[0] == "hello" {
		bot.Privmsg(channel, "world!")
	}
}
