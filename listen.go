package main

import (
	"log"
	"strings"

	"github.com/mferrera/go-ircevent"
	"github.com/mferrera/gobot/catfact"
	"github.com/mferrera/gobot/common"
	"github.com/mferrera/gobot/ohayou"
)

// check if string looks like an Irc channel
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
func listen(b common.Bot) {
	b.Irc.AddCallback("PRIVMSG", func(e *irc.Event) {
		admin := b.Admins[e.Nick] == e.Host
		word := strings.Split(e.Message(), " ")

		emit := common.EmitMsg{b.CommandPrefix, // the prefix
			word[0],        // the command word
			word,           // the entire message split into slice
			e.Arguments[0], // the channel (or PM) the message came from
			&b.Channels,
			e.Nick, // who typed the message
			admin,  // if the message was sent by an admin
			b.Irc,
			b.Irc.Privmsg,
			b.IgnoreList}

		// default admin commands
		go adminCommands(emit)

		// TODO: hostname/wildcard ignores not implemented yet
		if _, ok := b.IgnoreList[e.Nick]; ok {
			log.Println("Ignored message from", e.Nick)
		} else {
			go catfact.Run(emit)
			go ohayou.Run(emit)
		}
	})
}
