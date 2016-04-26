package main

import (
	"log"
	"strings"

	"github.com/mferrera/go-ircevent"
	"github.com/mferrera/gobot/catfact"
	"github.com/mferrera/gobot/ohayou"
)

// check if string looks like an irc channel
func isChannel(s string) bool {
	return strings.Index(s, "#") == 0
}

// check if message has more than one word
func hasArgs(a []string) bool {
	if len(a) > 1 {
		return true
	}

	return false
}

// acts as event emitter to all plugins
// for messages at least
func listen(bot *Bot) {
	bot.irc.AddCallback("PRIVMSG", func(e *irc.Event) {
		// split event.Message to array
		word := strings.Split(e.Message(), " ")

		// the command possibly being issued
		cmd := word[0]

		// prefix from conf
		p := bot.CommandPrefix
		channel := e.Arguments[0]
		nick := e.Nick

		// true or false if e.Nick@e.Host
		// is admin@adminhost
		admin := bot.Admins[e.Nick] == e.Host

		// default admin commands
		go bot.adminCommands(p, cmd, nick, channel, word, admin)

		// hostname/wildcard ignores not implemented yet
		if _, ok := bot.IgnoreList[e.Nick]; ok {
			log.Println("Ignored message from", e.Nick)
		} else {
			go catfact.Run(bot.irc, p, cmd, channel, word, admin)
			go ohayou.Run(bot.irc, p, cmd, channel, word, admin)
		}
	})
}
