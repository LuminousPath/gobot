package ohayou

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/mferrera/go-ircevent"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	eventsStarted bool   // if init() has started the events
	p             string // command prefix passed from main bot
	lowNick       string // nick lowercased
	argOne        string // word[1] lowercased
	argTwo        string // word[2] lowercased
	isPM          bool

	// for registering
	pin          int
	isRegistered bool
	getPin       = make(chan SubmitPin)

	// following set in newOhayou()
	typeResponse   string
	ohayous        int // new ohayous
	itemOhayous    int // extra ohayous given to user from items
	itemMultiplier int // any item multipliers
	totalOhayous   int // added all up

	inv        string   // used for inv command
	extra      string   // used in useItem
	err        error    // err var used everywhere for logging
	itemsInCat []string // slice used to return items in a category
	itemCats   []string // slice that holds all item categories
	top        []UserOhayous
	top5       string
	save       bson.M         // bson object that maps the "json" we save in DB queries
	t          time.Time      // time used everywhere
	est        *time.Location // timezone -- set in init

	adj = [11]string{"Great", "Superb", "Fantastic", "Amazing", "Marvelous",
		"Stunning", "Splendid", "Exquisite", "Impressive", "Outstanding", "Wonderful"}

	// DB "global" session -- all session are copied from this
	session *mgo.Session

	// user and item vars
	USER *User
	ITEM *Item

	// irc stuff
	b     *irc.Connection
	chans []string
	say   func(string, string)
)

func init() {
	// set timezone
	setEst, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic("Could not load TimeZone file")
	}
	est = setEst

	// set up DB session
	session, err = mgo.Dial(dbAddress)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	// fill up the slice of category names
	listCategories()
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

func isPin(pn string) bool {
	if len(pn) != 4 {
		return false
	}
	_, err = strconv.Atoi(pn)
	if err != nil {
		return false
	} else {
		return true
	}
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
		itemMultiplier = 1
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
		USER.saveOhayous(totalOhayous)
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

func Run(bot *irc.Connection, pre, cmd, channel, nick string, chnls, word []string) {
	b = bot
	say = b.Privmsg
	p = pre
	chans = chnls
	pin = 0
	lowNick = strings.ToLower(nick)
	if len(word) > 1 {
		argOne = strings.ToLower(word[1])
	}
	if len(word) > 2 {
		argTwo = strings.ToLower(word[2])
	}
	if channel[:1] == "#" {
		isPM = false
	} else {
		isPM = true
	}

	// check if events have started, if not, start them and set it so
	if !eventsStarted {
		startEvents()
		eventsStarted = true
	}

	// main command to acquire new ohayous
	if cmd == p+"ohayou" && !hasArgs(word) && !isPM {
		say(channel, newOhayou(nick))
	}

	// respond to channel with how many ohayous X has
	if cmd == p+"ohayou" && hasArgs(word) && !isPM {
		if getUser(argOne) {
			say(channel, fmt.Sprintf("%s has %d ohayous.", word[1], USER.Ohayous))
		} else {
			say(channel, argOne+" hasn't ohayoued yet!")
		}
	}

	if cmd == p+"buy" && !hasArgs(word) && !isPM {
		say(channel, "Usage: "+p+"buy <item> will buy you one <item>."+
			p+"buy <item> 3 will buy you 3 of <item>, if you can afford it.")
	} else if cmd == p+"buy" && hasArgs(word) && !isPM {
		// if a purchase quantity is given
		if len(word) > 2 {
			// try to convert it to an integer
			amt, err := strconv.Atoi(argTwo)
			if err != nil {
				say(channel, "You didn't give a valid quantity. Usage: "+p+
					"buy <item> will buy you one <item>. "+p+"buy <item>"+
					" 3 will buy you 3 of <item>, if you can afford it.")
			} else {
				say(channel, buyItem(lowNick, channel, argOne, amt))
			}
		} else {
			say(channel, buyItem(lowNick, channel, argOne, 1))
		}
	}

	// just shows how to use .items and lists item categories
	if cmd == p+"items" && !hasArgs(word) && !isPM {
		say(channel, "Type "+p+"item <category> to get a list of items by "+
			"category. Categories: "+strings.Join(append(itemCats), ", "))
	}

	// PMs all items in a category
	if cmd == p+"items" && hasArgs(word) {
		itemsInCat = getCategory(argOne)

		for _, itm := range itemsInCat {
			say(nick, itm)
		}
	}

	// returns information about an item
	if cmd == p+"item" && hasArgs(word) && !isPM {
		if getItem(argOne) {
			if ITEM.Purchase {
				say(channel, fmt.Sprintf("%s: %s - Price: %d ohayous",
					ITEM.Name, ITEM.Desc, ITEM.Price))
			} else {
				say(channel, fmt.Sprintf("%s: %s. Cannot be purchased.",
					ITEM.Name, ITEM.Desc))
			}
		} else {
			say(channel, "I don't carry that.")
		}
	}

	if cmd == p+"use" && !hasArgs(word) && !isPM {
		say(channel, "Type "+p+"use <item> to use an item. Type "+p+"inventory to "+
			"see what items you have, or "+p+"items to see what items you can "+
			p+"buy.")
	} else if cmd == p+"use" && hasArgs(word) && !isPM {
		if len(word) > 2 {
			say(channel, useItem(lowNick, nick, argOne, argTwo))
		} else {
			say(channel, useItem(lowNick, nick, argOne, "somebody"))
		}
	}

	// respond to nick with their items and quantity of each item
	if cmd == p+"inventory" {
		inv = "You have: "
		if !getUser(lowNick) {
			say(channel, "You haven't ohayoued yet! Type "+p+"ohayou to "+
				"get your first ration.")
		} else if len(USER.Items) > 0 {
			for itm, amt := range USER.Items {
				if amt == 0 {
					continue
				} else if amt > 1 {
					inv += fmt.Sprintf("%d %ss, ", amt, itm)
				} else {
					inv += fmt.Sprintf("%d %s, ", amt, itm)
				}
			}
			say(nick, inv[:len(inv)-2])
		} else {
			say(nick, "You don't have any items yet. Keep saving!")
		}
	}

	// say top 5 most ohayous at present
	if cmd == p+"top" && !isPM {
		getTop()
		top5 = "Top 5 Ohayouers: "
		for i := range top {
			top5 += fmt.Sprintf("%s: %d, ", top[i].Username, top[i].Ohayous)
		}
		say(channel, top5[:len(top5)-2])
	}

	if cmd == p+"register" {
		if !hasArgs(word) {
			say(nick, "Registering allows you to protect your ohayou assets. "+
				"Your nick must be registered to do so, and it will require "+
				"you to enter a pin number of your choosing for that command "+
				"to execute.")
			say(nick, "Type "+p+"register <pin> to register. The pin must be a "+
				"four digit number. DON'T USE YOUR REAL BANK PIN THOUGH "+
				"IDIOT. And remember to do it in PM!")
			say(nick, "Example: "+p+"register 1234")
		} else if hasArgs(word) && isPM {
			if len(argOne) != 4 {
				say(nick, "Your pin must be a four digit number. Example: "+
					p+"register 1234")
				return
			}

			// try to convert it to an integer
			pin, err = strconv.Atoi(argOne)
			if err != nil {
				say(nick, "Your pin must be a four digit number. Example: "+
					p+"register 1234")
			} else {
				doRegister(nick, pin)
			}

		}
	}

	if isPin(cmd) && isPM {
		pin, err = strconv.Atoi(cmd)
		if err != nil {
			return
		} else {
			getPin <- SubmitPin{lowNick, pin}
		}
	}
}
