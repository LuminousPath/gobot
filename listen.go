package main

import (
	"github.com/mferrera/gobot/hello"
	"github.com/thoj/go-ircevent"
	"log"
	"strings"
)

func isChannel(s string) bool {
	return strings.Index(s, "#") == 0
}

func hasArgs(a []string) bool {
	if len(a) > 1 {
		return true
	}

	return false
}

func listen(bot *Bot) {
	bot.irc.AddCallback("PRIVMSG", func(event *irc.Event) {
		// split event.Message to array
		word := strings.Split(event.Message(), " ")

		// the command possibly being issued
		cmd := word[0]
		// prefix from conf
		p := bot.CommandPrefix
		channel := event.Arguments[0]
		// true or false if event.Nick@event.Host
		// is admin@adminhost
		admin := bot.Admins[event.Nick] == event.Host

		// default admin commands
		go bot.adminCommands(p, cmd, channel, word, admin)

		// will be ignorelist soon
		if event.Nick == bot.irc.GetNick() {
			log.Println("Ignored message from", event.Nick)
		} else {
			go hello.Run(&bot.irc, p, cmd, channel, word, admin)
			go nude.Run(&bot.irc, p, cmd, channel, word, admin)
		}
	})
}
