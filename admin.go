package main

import (
	"strings"

	"github.com/mferrera/gobot/common"
)

func adminCommands(m common.EmitMsg) {
	if !m.Admin {
		return
	}

	say := m.Irc.Privmsg
	join := m.Irc.Join
	part := m.Irc.Part
	notice := m.Irc.Notice
	action := m.Irc.Action
	kick := m.Irc.Kick
	channels := *m.Channels

	// say #m.Channel hello
	// say hello -- says into m.Channel where command issued
	if m.Cmd == m.P+"say" && hasArgs(m.Word) {
		if isChannel(m.Word[1]) {
			say(m.Word[1], strings.Join(m.Word[2:], " "))
		} else {
			say(m.Channel, strings.Join(m.Word[1:], " "))
		}
	}

	// pm <user> hello
	if m.Cmd == m.P+"pm" && hasArgs(m.Word) {
		say(m.Word[1], strings.Join(m.Word[1:], " "))
	}

	// join #channel
	if m.Cmd == m.P+"join" && hasArgs(m.Word) {
		if isChannel(m.Word[1]) {
			join(m.Word[1])
			*m.Channels = append(*m.Channels, m.Word[1])
		}
	}

	if m.Cmd == m.P+"channels" {
		say(m.Channel, strings.Join(*m.Channels, " "))
	}

	// part -- parts channel where command issued
	// part #channel
	if m.Cmd == m.P+"part" {
		if !hasArgs(m.Word) {
			part(m.Channel)
			for i, c := range *m.Channels {
				if c == m.Channel {
					*m.Channels = append(channels[:i],
						channels[i+1:]...)
				}
			}
		} else if isChannel(m.Word[1]) {
			part(m.Word[1])
			for i, c := range *m.Channels {
				if c == m.Word[1] {
					*m.Channels = append(channels[:i],
						channels[i+1:]...)
				}
			}
		}
	}

	// notice <user> hello
	// notice #channel hello
	if m.Cmd == m.P+"notice" && hasArgs(m.Word) {
		if isChannel(m.Word[1]) {
			notice(m.Word[1], strings.Join(m.Word[2:], " "))
		} else {
			notice(m.Word[1], strings.Join(m.Word[2:], " "))
		}
	}

	// me hello
	// me #channel hello
	if m.Cmd == m.P+"me" && hasArgs(m.Word) {
		if !isChannel(m.Word[1]) {
			action(m.Channel, strings.Join(m.Word[1:], " "))
		} else if isChannel(m.Word[1]) {
			action(m.Word[1], strings.Join(m.Word[2:], " "))
		}
	}

	// kick <user> <reason>
	// kick #m.Channel <user> <reason> -- neither reason is required
	if m.Cmd == m.P+"kick" && hasArgs(m.Word) {
		if !isChannel(m.Word[1]) {
			kick(m.Word[1], m.Channel, strings.Join(m.Word[2:], " "))
		} else if isChannel(m.Word[1]) {
			kick(m.Word[2], m.Word[1], strings.Join(m.Word[3:], " "))
		}
	}

	// ignore <user> <reason>
	if m.Cmd == m.P+"ignore" {
		if len(m.Word) == 2 {
			m.IgnoreList[m.Word[1]] = "No reason given."

			say(m.Channel, "Added "+m.Word[1]+" to ignore list. "+
				"Reason: No reason given.")
		} else if len(m.Word) > 2 {
			reason := strings.Join(m.Word[2:], " ")
			m.IgnoreList[m.Word[1]] = reason

			say(m.Channel, "Added "+m.Word[1]+" to ignore list. "+
				"Reason: "+reason)
		}
	}

	// ignorelist
	if m.Cmd == m.P+"ignorelist" {
		var iglist string

		for n, r := range m.IgnoreList {
			iglist += "[ " + n + ", Reason: " + r + "] "
		}

		say(m.Nick, iglist)
	}

	// unignore <user>
	if m.Cmd == m.P+"unignore" && hasArgs(m.Word) {
		if _, ok := m.IgnoreList[m.Word[1]]; ok {
			delete(m.IgnoreList, m.Word[1])
			say(m.Channel, "Unignored "+m.Word[1]+".")
		} else {
			say(m.Channel, m.Word[1]+" is not ignored.")
		}
	}
}
