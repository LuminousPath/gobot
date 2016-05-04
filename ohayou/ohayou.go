package ohayou

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

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

// main function that distributes ohayous
func newOhayou(nick string) string {
	adj := [11]string{"Great", "Superb", "Fantastic", "Amazing", "Marvelous", "Stunning", "Splendid", "Exquisite", "Impressive", "Outstanding", "Wonderful"}
	ohayous := randNum(0, 7)

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

	// get their data
	data := getUser(strings.ToLower(nick))

	t := time.Now()
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("err: ", err.Error())
	}

	// dont allow ohayou if they have ohayou'd today
	if data.Last.Format("20060102") >= t.In(est).Format("20060102") {
		return fmt.Sprintf("You already got your ohayou ration today, "+
			"%s. If you'd like to ohayou again, you can purchase another "+
			".ohayou for $5.00 USD.", nick)
	}

	if data.TimesOhayoued == 0 {
		newUser(strings.ToLower(nick), ohayous)
		return fmt.Sprintf("Congratulations on your first ohayou %s!!! "+
			"%s Type .ohayouhelp if you don't know what this is.",
			nick, typeResponse)
	} else {
		totalOhayous := saveOhayous(data, ohayous)
		return fmt.Sprintf("%s ohayou %s!!! %s You have %d ohayous.",
			adj[randNum(0, 11)], nick, typeResponse, totalOhayous)
	}
}

func Run(bot *irc.Connection, p, cmd, channel, nick string, word []string, admin bool) {
	say := bot.Privmsg

	// main command to acquire new ohayous
	if cmd == p+"ohayou" && !hasArgs(word) {
		say(channel, newOhayou(nick))
	}

	// respond to channel with how many ohayous X has
	if cmd == p+"ohayou" && hasArgs(word) {
		data := getUser(strings.ToLower(word[1]))
		if data.Username != "" {
			say(channel, fmt.Sprintf("%s has %d ohayous.", word[1], data.Ohayous))
		} else {
			say(channel, word[1]+" hasn't ohayoued yet!")

		}
	}

	// respond to nick with their items and quantity of each item
	if cmd == p+"inventory" {
		data := getUser(strings.ToLower(nick))
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
}
