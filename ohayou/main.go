package ohayou

import (
	"math/rand"
	"time"

	"github.com/mferrera/gobot/common"
	"gopkg.in/mgo.v2"
)

var (
	p        string         // command p+refix p+assed from main bot
	itemCtgs []string       // slice that holds all item categories
	est      *time.Location // timezone -- set in init

	identified    = make(map[string]bool) // map of identified users
	watchingNicks bool

	// irc stuff
	say   func(string, string)
	chans []string

	// persitent db session
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

func isPM(m common.EmitMsg) bool {
	if m.Channel[:1] == "#" {
		return false
	}
	return true
}

func Run(m common.EmitMsg) {
	say = m.Say
	p = m.P
	chans = *m.Channels

	switch m.Cmd {
	case p + "changelog":
		say(m.Channel, "Latest changelog: http://pastebin.com/raw/b99DQ0KZ")
	case p + "help":
		cmd_help(m)
	case p + "ohayou":
		cmd_ohayou(m)
	case p + "buy":
		cmd_buy(m)
	case p + "equip":
		cmd_equip(m)
	case p + "unequip":
		cmd_unequip(m)
	case p + "items":
		cmd_items(m)
	case p + "item":
		cmd_item(m)
	case p + "use":
		cmd_use(m)
	case p + "steal":
		cmd_steal(m)
	case p + "deposit":
		cmd_deposit(m)
	case p + "withdraw":
		cmd_withdraw(m)
	case p + "stats":
		cmd_stats(m)
	case p + "inventory":
		cmd_inventory(m)
	case p + "register":
		cmd_register(m)
	case p + "identify":
		cmd_identify(m)
	case p + "quarry":
		cmd_quarry(m)
	}
}
