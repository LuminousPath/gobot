package ohayou

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/mferrera/go-ircevent"
)

func hasArgs(a []string) bool {
	if len(a) > 1 {
		return true
	}

	return false
}

func randNum(min, max int) int {
	return min + rand.Intn(max-min)
}

func newOhayou(nick string) string {
	adj := [11]string{"Great", "Superb", "Fantastic", "Amazing", "Marvelous", "Stunning", "Splendid", "Exquisite", "Impressive", "Outstanding", "Wonderful"}
	ohayous := randNum(0, 6)

	typeResponse := ""
	switch ohayous {
	case 0:
		typeResponse = "But not good enough. You get 0 ohayous."
	case 1:
		typeResponse = "You get 1 ohayou."
	case 6:
		typeResponse = "Wow! You get 6 ohayous!"
	default:
		typeResponse = "You get " + strconv.Itoa(ohayous) + " ohayous!"
	}

	return fmt.Sprintf("%s ohayou %s!!! %s You get %d ohayous.",
		adj[randNum(0, 11)], nick, typeResponse, ohayous)
}

func Run(bot *irc.Connection, p, cmd, channel, nick string, word []string, admin bool) {
	say := bot.Privmsg

	// main command to acquire new ohayous
	if cmd == p+"ohayou" && !hasArgs(word) {
		say(channel, newOhayou(nick))
	}

	// respond to channel with how many ohayous X has
	if cmd == p+"ohayou" && hasArgs(word) {
		data := getUser(word[1])
		say(channel, fmt.Sprintf("%s has %d ohayous.", data.Username, data.Ohayous))
	}

	// respond to nick with their items and quantity of each item
	if cmd == p+"inventory" {
		data := getUser(nick)
		inv := "You have "

		for item, qty := range data.Items {
			if qty > 1 {
				inv += strconv.Itoa(qty) + " " + item + "s "
			} else {
				inv += strconv.Itoa(qty) + " " + item + " "
			}
		}

		say(nick, inv)
	}

	return
}
