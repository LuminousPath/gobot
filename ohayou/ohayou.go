package ohayou

import (
	"github.com/mferrera/go-ircevent"
)

func hasArgs(a []string) bool {
	if len(a) > 1 {
		return true
	}

	return false
}

func Run(bot *irc.Connection, p, cmd, channel string, word []string, admin bool) {
	say := bot.Privmsg

	if cmd == p+"ohayou" && hasArgs(word) {
		say(channel, getOhayous(word[1]))
	}

	return
}
