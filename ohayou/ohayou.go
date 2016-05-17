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

// globalize things that will be used repeatedly throughout the package
var (
	p       string // command prefix passed from main bot
	lowNick string // nick lowercased
	argOne  string // word[1] lowercased
	argTwo  string // word[2] lowercased
	isPM    bool

	// for registering
	pin    int
	getPin = make(chan SubmitPin)

	inv        string        // used for inv command
	err        error         // err var used everywhere for logging
	itemsInCat []string      // slice used to return items in a category
	itemCats   []string      // slice that holds all item categories
	top        []UserOhayous // holds users+ohayous to be iterated over
	top5       string        // iterated over user+ohayou struct and concatenated
	save       bson.M        // bson object that maps the "json" we save in DB queries
	t          time.Time     // time used everywhere
	last       time.Time
	now        time.Time
	est        *time.Location // timezone -- set in init

	// irc stuff
	b       *irc.Connection
	chans   []string
	say     func(string, string)
	curChan string
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
	go listCategories()
	go startEvents()  // start all special events
	go fillFortunes() // fill fortunes var
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

func Run(bot *irc.Connection, pre, cmd, channel, nick string, chnls, word []string) {
	b = bot
	say = b.Privmsg
	p = pre
	curChan = channel
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

	// latest changelog
	if cmd == p+"changelog" {
		say(channel, "Latest changelog: http://pastebin.com/LANmT0Ww")
	}

	// displays some help
	if cmd == p+"help" && hasArgs(word) {
		if argOne == "ohayou" {
			say(channel, "An ohayou game. Acquire ohayous and purchase things "+
				"with them. Some items have special functions. Commands: "+p+
				"ohayou, "+p+"buy, "+p+"item, "+p+"items, "+p+"use, "+p+
				"inventory, "+p+"register, "+p+"changelog")
		}
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
		say(channel, "Usage: "+p+"buy <item> will buy you one <item>. "+
			p+"buy <item> 3 will buy you 3 of <item>, if you can afford it.")
	} else if cmd == p+"buy" && hasArgs(word) && !isPM {
		if argOne == "ohayou" && getUser(lowNick) {
			say(channel, fmt.Sprintf("You purchased %d ohayous for %d ohayous. "+
				"You have %d ohayous left.",
				USER.Ohayous, USER.Ohayous, USER.Ohayous))
			return
		}
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
		say(channel, "Type "+p+"items <category> to get a list of items by "+
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

	if cmd == p+"steal" && hasArgs(word) {
		if getUser(lowNick) {
			if PutUser(argOne, &stealVictim) {
				go USER.StealFrom(stealVictim, channel, nick, word[1])
			} else {
				say(channel, "You can't steal from "+word[1]+" because "+
					word[1]+"has never ohayou'd!")
			}
		} else {
			say(channel, "You can't do that because you haven't ohayou'd yet! "+
				"Type "+p+"ohayou to get your first ration.")
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
				doRegister(lowNick, pin)
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
