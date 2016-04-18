package main

import (
	"strings"
)

func (bot *Bot) adminCommands(p, cmd, channel string, word []string, admin bool) {
	if !admin {
		return
	}

	say := bot.irc.Privmsg
	join := bot.irc.Join
	part := bot.irc.Part
	notice := bot.irc.Notice
	action := bot.irc.Action
	kick := bot.irc.Kick

	if cmd == p+"say" {
		if hasArgs(word) && isChannel(word[1]) {
			say(word[1], strings.Join(word[2:], " "))
		} else if hasArgs(word) {
			say(channel, strings.Join(word[1:], " "))
		}
	}

	if cmd == p+"pm" {
		if hasArgs(word) {
			say(word[1], strings.Join(word[1:], " "))
		}
	}

	if cmd == p+"join" {
		if hasArgs(word) && isChannel(word[1]) {
			join(word[1])
		}
	}

	if cmd == p+"part" {
		if !hasArgs(word) {
			part(channel)
		} else if isChannel(word[1]) {
			part(word[1])
		}
	}

	if cmd == p+"notice" {
		if hasArgs(word) && isChannel(word[1]) {
			notice(word[1], strings.Join(word[2:], " "))
		} else if hasArgs(word) {
			notice(word[1], strings.Join(word[2:], " "))
		}
	}

	if cmd == p+"me" {
		if hasArgs(word) && !isChannel(word[1]) {
			action(channel, strings.Join(word[1:], " "))
		} else if hasArgs(word) && isChannel(word[1]) {
			action(word[1], strings.Join(word[2:], " "))
		}
	}

	if cmd == p+"kick" {
		if hasArgs(word) && !isChannel(word[1]) {
			kick(word[1], channel, strings.Join(word[2:], " "))
		} else if hasArgs(word) && isChannel(word[1]) {
			kick(word[2], word[1], strings.Join(word[3:], " "))
		}
	}

	return
}
