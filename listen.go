package main

import (
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
		if event.Nick == bot.irc.GetNick() {
			log.Println("Ignored message from", event.Nick)
		} else {
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

			// simplified functions
			say := bot.irc.Privmsg
			//sayf := bot.irc.Privmsgf
			join := bot.irc.Join
			part := bot.irc.Part
			notice := bot.irc.Notice
			//noticef := bot.irc.Noticef
			action := bot.irc.Action
			//actionf := bot.irc.Actionf
			kick := bot.irc.Kick
			//multikick := bot.irc.MultiKick
			//getnick := bot.irc.GetNick
			//whois := bot.irc.Whois
			//mode := bot.irc.Mode

			// default admin commands
			if cmd == p+"say" && admin {
				if hasArgs(word) && isChannel(word[1]) {
					say(word[1], strings.Join(word[2:], " "))
				} else if hasArgs(word) {
					say(channel, strings.Join(word[1:], " "))
				}
			}

			if cmd == p+"pm" && admin {
				if hasArgs(word) {
					say(word[1], strings.Join(word[1:], " "))
				}
			}

			if cmd == p+"join" && admin {
				if hasArgs(word) && isChannel(word[1]) {
					join(word[1])
				}
			}

			if cmd == p+"part" && admin {
				if !hasArgs(word) {
					part(channel)
				} else if isChannel(word[1]) {
					part(word[1])
				}
			}

			if cmd == p+"notice" && admin {
				if hasArgs(word) && isChannel(word[1]) {
					notice(word[1], strings.Join(word[2:], " "))
				} else if hasArgs(word) {
					notice(word[1], strings.Join(word[2:], " "))
				}
			}

			if cmd == p+"me" && admin {
				if hasArgs(word) && !isChannel(word[1]) {
					action(channel, strings.Join(word[1:], " "))
				} else if hasArgs(word) && isChannel(word[1]) {
					action(word[1], strings.Join(word[2:], " "))
				}
			}

			if cmd == p+"kick" && admin {
				if hasArgs(word) && !isChannel(word[1]) {
					kick(word[1], channel, strings.Join(word[2:], " "))
				} else if hasArgs(word) && isChannel(word[1]) {
					kick(word[2], word[1], strings.Join(word[3:], " "))
				}
			}
		}
	})
}
