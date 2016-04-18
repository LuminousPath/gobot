package main

import (
	"strings"
)

func (bot *Bot) adminCommands(p, cmd, nick, channel string, word []string, admin bool) {
	if !admin {
		return
	}

	say := bot.irc.Privmsg
	join := bot.irc.Join
	part := bot.irc.Part
	notice := bot.irc.Notice
	action := bot.irc.Action
	kick := bot.irc.Kick

	// say #channel hello
	// say hello -- says into channel where command issued
	if cmd == p+"say" && hasArgs(word) {
		if isChannel(word[1]) {
			say(word[1], strings.Join(word[2:], " "))
		} else {
			say(channel, strings.Join(word[1:], " "))
		}
	}

	// pm <user> hello
	if cmd == p+"pm" && hasArgs(word) {
		say(word[1], strings.Join(word[1:], " "))
	}

	// join #channel
	if cmd == p+"join" && hasArgs(word) {
		if isChannel(word[1]) {
			join(word[1])
		}
	}

	// part -- parts channel where command issued
	// part #channel
	if cmd == p+"part" {
		if !hasArgs(word) {
			part(channel)
		} else if isChannel(word[1]) {
			part(word[1])
		}
	}

	// notice <user> hello
	// notice #channel hello
	if cmd == p+"notice" && hasArgs(word) {
		if isChannel(word[1]) {
			notice(word[1], strings.Join(word[2:], " "))
		} else {
			notice(word[1], strings.Join(word[2:], " "))
		}
	}

	// me hello
	// me #channel hello
	if cmd == p+"me" && hasArgs(word) {
		if !isChannel(word[1]) {
			action(channel, strings.Join(word[1:], " "))
		} else if isChannel(word[1]) {
			action(word[1], strings.Join(word[2:], " "))
		}
	}

	// kick <user> <reason>
	// kick #channel <user> <reason> -- neither reason is required
	if cmd == p+"kick" && hasArgs(word) {
		if !isChannel(word[1]) {
			kick(word[1], channel, strings.Join(word[2:], " "))
		} else if isChannel(word[1]) {
			kick(word[2], word[1], strings.Join(word[3:], " "))
		}
	}

	// ignore <user> <reason>
	if cmd == p+"ignore" {
		if len(word) == 2 {
			bot.IgnoreList[word[1]] = "No reason given."

			say(channel, "Added "+word[1]+" to ignore list. "+
				"Reason: No reason given.")
		} else if len(word) > 2 {
			reason := strings.Join(word[2:], " ")
			bot.IgnoreList[word[1]] = reason

			say(channel, "Added "+word[1]+" to ignore list. "+
				"Reason: "+reason)
		}
	}

	// ignorelist
	if cmd == p+"ignorelist" {
		var iglist string

		for n, r := range bot.IgnoreList {
			iglist += "[ " + n + ", Reason: " + r + "] "
		}

		say(nick, iglist)
	}

	// unignore <user>
	if cmd == p+"unignore" && hasArgs(word) {
		if _, ok := bot.IgnoreList[word[1]]; ok {
			delete(bot.IgnoreList, word[1])
			say(channel, "Unignored "+word[1]+".")
		} else {
			say(channel, word[1]+" is not ignored.")
		}
	}

	return
}
