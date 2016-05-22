package ohayou

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/mferrera/gobot/common"
	"gopkg.in/mgo.v2"
)

var (
	p string // command p+refix p+assed from main bot

	// for registering
	sendPin = make(chan SubmitPin)

	itemCtgs []string       // slice that holds all item categories
	est      *time.Location // timezone -- set in init

	// irc stuff
	say   func(string, string)
	chans []string

	session *mgo.Session
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

	go setCategories() // fill var with current item categories
	go fillFortunes()  // p+ut all fortunes in var
	go startEvents()   // start all special events
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

func randBigNum(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func isPin(pin string) bool {
	if len(pin) != 4 {
		return false
	}
	_, err := strconv.Atoi(pin)
	if err != nil {
		return false
	} else {
		return true
	}
}

func Run(m common.EmitMsg) {
	say = m.Say
	p = m.P
	chans = *m.Channels
	lowNick := strings.ToLower(m.Nick)
	var argOne, argTwo string
	var isPM bool

	// lowercase first and second arguments
	// since they are often items/m.Nicks
	if len(m.Word) > 1 {
		argOne = strings.ToLower(m.Word[1])
	}
	if len(m.Word) > 2 {
		argTwo = strings.ToLower(m.Word[2])
	}

	// check if this is a PM
	if m.Channel[:1] == "#" {
		isPM = false
	} else {
		isPM = true
	}

	// latest changelog
	if m.Cmd == p+"changelog" {
		say(m.Channel, "Latest changelog: http://pastebin.com/raw/b2dh6Prt")
	}

	// displays some help
	if m.Cmd == p+"help" && argOne == "ohayou" {
		say(m.Channel, "An ohayou game. Acquire ohayous and purchase things "+
			"with them. Some items have special functions. Commands: "+p+
			"ohayou, "+p+"buy, "+p+"item, "+p+"items, "+p+"use, "+p+
			"inventory, "+p+"register, "+p+"changelog")
	}

	// main command to acquire new ohayous
	if m.Cmd == p+"ohayou" && !hasArgs(m.Word) && !isPM {
		say(m.Channel, Ohayou(lowNick))
	}

	if m.Cmd == p+"ohayous" {
		user, ok := GetUser(lowNick)
		if ok {
			if user.Vault.Installed {
				say(m.Nick, fmt.Sprintf("You have %d ohayous on hand and %d "+
					"ohayous in your Level %d vault. Your vault's "+
					"capacity is %d ohayous.",
					user.Ohayous, user.Vault.Ohayous, user.Vault.Level+1,
					int(math.Pow(10, 2+float64(user.Vault.Level)))))
			} else {
				say(m.Nick, fmt.Sprintf("You have %d ohayous.", user.Ohayous))
			}
		} else {
			say(m.Channel, "You haven't ohayou'd yet! Type "+p+"ohayou "+
				"to get your first ration.")
		}
	}

	if m.Cmd == p+"buy" && !hasArgs(m.Word) && !isPM {
		say(m.Channel, "Usage: "+p+"buy <item> will buy you one <item>. "+
			p+"buy <item> 3 will buy you 3 of <item>, if you can afford it.")
	} else if m.Cmd == p+"buy" && hasArgs(m.Word) && !isPM {
		user, ok := GetUser(lowNick)
		if ok {
			// just for fun
			if argOne == "ohayou" {
				say(m.Channel, fmt.Sprintf("You purchased %d ohayous "+
					"for %d ohayous. You have %d ohayous left.",
					user.Ohayous, user.Ohayous, user.Ohayous))
				return
			}
			// if a p+urchase quantity is given
			if len(m.Word) > 2 {
				// try to convert it to an integer
				amt, err := strconv.Atoi(argTwo)
				if err != nil {
					say(m.Channel, "You didn't give a valid quantity. "+
						"Usage: "+p+"buy <item> will buy you one "+
						"<item>. "+p+"buy <item>"+" 3 will buy you "+
						"3 of <item>, if you can afford it.")
				} else {
					say(m.Channel, user.Buy(m.Channel, argOne, amt))
				}
			} else {
				say(m.Channel, user.Buy(m.Channel, argOne, 1))
			}
		} else {
			say(m.Channel, "You haven't ohayoued yet! Type "+p+"ohayou to get"+
				" your first ration.")
		}
	}

	if m.Cmd == p+"equip" {
		user, ok := GetUser(lowNick)
		if !ok {
			say(m.Channel, "You haven't ohayou'd yet! Type "+p+"ohayou to get "+
				"your first ration.")
			return
		}
		if !hasArgs(m.Word) {
			if isPM {
				say(m.Nick, "Type "+p+"equip <item> to equip <item>"+
					"You can only have one item equipped per slot, "+
					"unless otherwise noted.")
			} else {
				say(m.Channel, "Type "+p+"equip <item> to equip <item>."+
					"You can only have one item equipped per "+
					"slot, unless otherwise noted.")
			}
		} else {
			if isPM {
				say(m.Nick, user.Equip(argOne))
			} else {
				say(m.Channel, user.Equip(argOne))
			}
		}
	}

	if m.Cmd == p+"unequip" {
		user, ok := GetUser(lowNick)
		if !ok {
			say(m.Channel, "You haven't ohayou'd yet! Type "+p+"ohayou to get "+
				"your first ration.")
			return
		}
		if !hasArgs(m.Word) {
			if isPM {
				say(m.Nick, "Type "+p+"unequip <item> to unequip <item>")
			} else {
				say(m.Channel, "Type "+p+"unequip <item> to unequip <item>")
			}
		} else {
			if isPM {
				say(m.Nick, user.Unequip(argOne))
			} else {
				say(m.Channel, user.Unequip(argOne))
			}
		}
	}

	// just shows how to use .items and lists item categories
	if m.Cmd == p+"items" && !hasArgs(m.Word) && !isPM {
		say(m.Channel, "Type "+p+"items <category> to get a list of items by "+
			"category. Categories: "+strings.Join(append(itemCtgs), ", ")+".")
	}

	// PMs all items in a category
	if m.Cmd == p+"items" && hasArgs(m.Word) {
		itemsInCtg := ItemCategory(argOne)
		for _, itm := range itemsInCtg {
			say(m.Nick, itm)
		}
	}

	// returns information about an item
	if m.Cmd == p+"item" && hasArgs(m.Word) && !isPM {
		item, ok := GetItem(argOne)
		if ok {
			if item.Purchase {
				say(m.Channel, fmt.Sprintf("%s: %s - Price: %d ohayous.",
					item.Name, item.Desc, item.Price))
			} else {
				say(m.Channel, fmt.Sprintf("%s: %s. Cannot be purchased.",
					item.Name, item.Desc))
			}
		} else {
			say(m.Channel, "I don't carry that.")
		}
	}

	if m.Cmd == p+"use" && !hasArgs(m.Word) && !isPM {
		say(m.Channel, "Type "+p+"use <item> to use an item. Type "+p+
			"inventory to see what items you have, or "+p+"items to see what "+
			"items you can "+p+"buy.")
	} else if m.Cmd == p+"use" && hasArgs(m.Word) && !isPM {
		user, ok := GetUser(lowNick)
		if ok {
			if len(m.Word) > 2 {
				say(m.Channel, user.Use(m.Nick, argOne, argTwo))
			} else {
				say(m.Channel, user.Use(m.Nick, argOne, "somebody"))
			}
		} else {
			say(m.Channel, "You don't have any items because you haven't "+
				"ohayou'd yet! Get your first ration by typing "+p+"ohayou.")
		}
	}

	if m.Cmd == p+"steal" {
		if !hasArgs(m.Word) {
			say(m.Channel, "Attempts to steal from someone. Usage: "+p+
				"steal <user>. Has penalties if you are caught!")
		} else {
			user, ok := GetUser(lowNick)
			if ok {
				victim, alsoOk := GetUser(argOne)
				if alsoOk {
					go user.StealFrom(victim, m.Channel, m.Nick, m.Word[1])
				} else {
					say(m.Channel, "You can't steal from "+m.Word[1]+" "+
						"because "+m.Word[1]+" has never ohayou'd!")
				}
			} else {
				say(m.Channel, "You can't do that because you haven't "+
					"ohayou'd yet! Type "+p+"ohayou to get your first "+
					"ration.")
			}
		}
	}

	if m.Cmd == p+"deposit" {
		if !hasArgs(m.Word) {
			say(m.Channel, "Deposits ohayous to your vault. Usage: "+p+
				"deposit <num> -- deposits <num> ohayous. Your vault "+
				"can only be opened once per day due to its security "+
				"protocol.")
		} else {
			user, ok := GetUser(lowNick)
			if ok {
				amt, err := strconv.Atoi(argOne)
				if err != nil {
					say(m.Channel, "You didn't give a valid quantity. "+
						"Usage: "+p+"deposit <num> will deposit "+
						"<num> ohayous to your secured storage.")
				} else {
					if isPM {
						say(m.Nick, user.Deposit(amt))
					} else {
						say(m.Channel, user.Deposit(amt))
					}
				}
			}
		}
	}

	if m.Cmd == p+"withdraw" {
		if !hasArgs(m.Word) {
			say(m.Channel, "Withdraws ohayous from your vault. Usage: "+p+
				"withdraw <num> -- withdraws <num> ohayous. Your vault "+
				"can only be opened once per day due to its security "+
				"protocol.")
		} else {
			user, ok := GetUser(lowNick)
			if ok {
				amt, err := strconv.Atoi(argOne)
				if err != nil {
					say(m.Channel, "You didn't give a valid quantity. "+
						"Usage: "+p+"withdraw <num> will withdraw "+
						"<num> ohayous to your secured storage.")
				} else {
					if isPM {
						say(m.Nick, user.Withdraw(amt))
					} else {
						say(m.Channel, user.Withdraw(amt))
					}
				}
			}
		}
	}

	if m.Cmd == p+"stats" {
		user, ok := GetUser(lowNick)
		if ok {
			go user.Stats()
		} else {
			say(m.Channel, "You can't do that because you haven't ohayou'd "+
				"yet! Type "+p+"ohayou to get your first ration.")
		}
	}

	// respond to m.Nick with their items and quantity of each item
	if m.Cmd == p+"inventory" {
		user, ok := GetUser(lowNick)
		if !ok {
			say(m.Channel, "You haven't ohayoued yet! Type "+p+"ohayou to "+
				"get your first ration.")
		} else if len(user.Items) > 0 {
			inv := "You have: "
			if user.Vault.Installed {
				inv += fmt.Sprintf("a Level %d vault (%d/%d ohayous), ",
					user.Vault.Level+1, user.Vault.Ohayous,
					int(math.Pow(10, 2+float64(user.Vault.Level))))
			}
			for itm, amt := range user.Items {
				if amt == 0 {
					continue
				} else if amt > 1 {
					inv += fmt.Sprintf("%d %ss, ", amt, itm)
				} else {
					inv += fmt.Sprintf("%d %s, ", amt, itm)
				}
			}
			say(m.Nick, inv[:len(inv)-2])
		} else {
			say(m.Nick, "You don't have any items yet. Keep saving!")
		}
	}

	if m.Cmd == p+"register" {
		if !hasArgs(m.Word) {
			say(m.Nick, "Registering allows you to protect your ohayou "+
				"assets. Your nick must be registered to do so, and it will"+
				"require you to enter a pin number of your choosing for that "+
				"command to execute.")
			say(m.Nick, "Type "+p+"register <pin> to register. The pin must "+
				"be a four digit number. DON'T USE YOUR REAL BANK PIN THOUGH "+
				"IDIOT. And remember to do it in PM!")
			say(m.Nick, "Example: "+p+"register 1234")
		} else if hasArgs(m.Word) && isPM {
			user, ok := GetUser(lowNick)
			if !ok {
				say(m.Nick, "Looks like you haven't ohayou'd yet. Type "+p+
					"ohayou in a channel I'm in to get your ration, and "+
					"then you can register.")
				return
			} else if len(argOne) != 4 {
				say(m.Nick, "Your pin must be a four digit number. "+
					"Example: "+p+"register 1234")
				return
			}

			// try to convert it to an integer
			pin, err := strconv.Atoi(argOne)
			if err != nil {
				say(m.Nick, "Your pin must be a four digit number. "+
					"Example: "+p+"register 1234")
			} else {
				user.Register(pin, m.Irc)
			}

		}
	}

	if isPin(m.Cmd) && isPM {
		pin, err := strconv.Atoi(m.Cmd)
		if err != nil {
			return
		} else {
			sendPin <- SubmitPin{lowNick, pin}
		}
	}
}
