package ohayou

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/mferrera/go-ircevent"
	"gopkg.in/mgo.v2"
)

// db consts
const (
	dbAddress string = "localhost"
	dbName    string = "ircbot"
	ohyCol    string = "ohayou"
	itemCol   string = "items"
)

// globalize things that will be used repeatedly throughout the package
var (
	p              string
	eventsStarted  bool
	argOne         string
	argTwo         string
	typeResponse   string
	ohayous        int
	itemOhayous    int
	itemMultiplier int = 1
	totalOhayous   int
	lowNick        string
	inv            string = "You have "
	extra          string
	err            error
	itemsInCat     []string
	t              time.Time
	est            *time.Location

	adj = [11]string{"Great", "Superb", "Fantastic", "Amazing", "Marvelous",
		"Stunning", "Splendid", "Exquisite", "Impressive", "Outstanding", "Wonderful"}

	// DB vars
	session   *mgo.Session
	dbInitErr error

	USER *User
	ITEM *Item
)

func init() {
	// set timezone
	setEst, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic("Could not load TimeZone file")
	}
	est = setEst

	// set up DB session
	session, dbInitErr = mgo.Dial(dbAddress)
	if dbInitErr != nil {
		panic(dbInitErr)
	}
	session.SetMode(mgo.Monotonic, true)
}

func clearGlobals() {
	inv = "You have "
	extra = ""
}

func hasArgs(a []string) bool {
	if len(a) > 1 {
		return true
	}
	return false
}

func randNum(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// main function that distributes ohayous
func newOhayou(nick string) string {
	ohayous = randNum(0, 6)

	switch ohayous {
	case 0:
		typeResponse = "But not good enough. You get 0 ohayous."
	case 1:
		typeResponse = "You get 1 ohayou."
	case 6:
		typeResponse = "Wow! You get 6 ohayous!"
	default:
		typeResponse = fmt.Sprintf("You get %d ohayous!", ohayous)
	}

	// get their data
	t = time.Now()

	// dont allow ohayou if they have ohayou'd today
	if !getUser(lowNick) {
		newUser(lowNick, ohayous)

		return "Congratulations on your first ohayou " + nick + "!!! " +
			typeResponse + " Type .ohayouhelp if you don't know what this is."
	} else if USER.Last.Format("20060102") >= t.In(est).Format("20060102") {
		return "You already got your ohayou ration today, " + nick
	} else {
		for itm, amt := range USER.Items {
			// check if user has item(s) that multiply another item
			if USER.ItemMultiply[itm] != 0 {
				itemMultiplier = USER.ItemMultiply[itm]
			}

			getItem(itm)
			itemOhayous += (ITEM.Add * amt) * itemMultiplier
		}

		totalOhayous = USER.Ohayous + ohayous + itemOhayous

		// store it
		saveOhayous(USER, totalOhayous)

		if itemOhayous == 0 {
			return fmt.Sprintf("%s ohayou %s!!! %s You have %d ohayous.",
				adj[randNum(0, 11)], nick, typeResponse, totalOhayous)
		} else {
			return fmt.Sprintf("%s ohayou %s!!! %s Your items increased "+
				"that to %d. You have %d ohayous.",
				adj[randNum(0, 11)], nick, typeResponse,
				ohayous+itemOhayous, totalOhayous)
		}
	}
}

func Run(bot *irc.Connection, prefix, cmd, channel, nick string, word []string, admin bool) {
	say := bot.Privmsg
	// make the prefix global
	p = prefix
	lowNick = strings.ToLower(nick)
	if len(word) > 1 {
		argOne = strings.ToLower(word[1])
	}
	if len(word) > 2 {
		argTwo = strings.ToLower(word[2])
	}

	// check if events have started, if not, start them and set it so
	if !eventsStarted {
		startEvents(bot)
		eventsStarted = true
	}

	// main command to acquire new ohayous
	if cmd == p+"ohayou" && !hasArgs(word) {
		say(channel, newOhayou(nick))
	}

	// respond to channel with how many ohayous X has
	if cmd == p+"ohayou" && hasArgs(word) {
		if getUser(argOne) {
			say(channel, fmt.Sprintf("%s has %d ohayous.", word[1], USER.Ohayous))
		} else {
			say(channel, argOne+" hasn't ohayoued yet!")
		}
	}

	if cmd == p+"buy" && !hasArgs(word) {
		say(channel, "Usage: "+p+"buy <item> will buy you one <item>."+
			p+"buy <item> 3 will buy you 3 of <item>, if you can afford it.")
	} else if cmd == p+"buy" && hasArgs(word) {
		// if a purchase quantity is given
		if len(word) > 2 {
			// try to convert it to an integer
			amt, err := strconv.Atoi(argTwo)
			if err != nil {
				say(channel, "You didn't give a valid quantity. Usage: "+p+
					"buy <item> will buy you one <item>. "+p+"buy <item>"+
					" 3 will buy you 3 of <item>, if you can afford it.")
			} else {
				say(channel, buyItem(lowNick, argOne, amt))
			}
		} else {
			say(channel, buyItem(lowNick, argOne, 1))
		}
	}

	// just shows how to use .items and lists item categories
	if cmd == p+"items" && !hasArgs(word) {
		say(channel, "Type "+p+"item <category> to get a list of items by "+
			"category. Categories: "+listCategories())
	}

	// PMs all items in a category
	if cmd == p+"items" && hasArgs(word) {
		itemsInCat = getCategory(argOne)

		for _, itm := range itemsInCat {
			say(nick, itm)
		}
	}

	// returns information about an item
	if cmd == p+"item" && hasArgs(word) {
		if getItem(argOne) {
			say(channel, fmt.Sprintf("%s: %s - Price: %d ohayous",
				ITEM.Name, ITEM.Desc, ITEM.Price))
		} else {
			say(channel, "I don't carry that.")
		}
	}

	if cmd == p+"use" && !hasArgs(word) {
		say(channel, "Type "+p+"use <item> to use an item. Type "+p+"inventory to "+
			"see what items you have, or "+p+"items to see what items you can "+
			p+"buy.")
	} else if cmd == p+"use" && hasArgs(word) {
		if len(word) > 2 {
			say(channel, useItem(lowNick, nick, argOne, argTwo))
		} else {
			say(channel, useItem(lowNick, nick, argOne, "somebody"))
		}
	}

	// respond to nick with their items and quantity of each item
	if cmd == p+"inventory" {
		getUser(lowNick)

		if USER.TimesOhayoued == 0 {
			say(channel, "You haven't ohayoued yet! Type .ohayou to "+
				"get your first ration.")
		} else if len(USER.Items) > 0 {
			for itm, amt := range USER.Items {
				if amt == 0 {
					continue
				} else if amt > 1 {
					inv += fmt.Sprintf("%d %ss ", amt, itm)
				} else {
					inv += fmt.Sprintf("%d %s ", amt, itm)
				}
			}

			say(nick, inv)
		} else {
			say(nick, "You don't have any items yet. Keep saving!")
		}
	}
	clearGlobals()
}
