package ohayou

import (
	"fmt"
	"time"
)

func buyItem(nick, channel, itm string, amt int) string {
	if !getUser(nick) {
		return "You haven't ohayoued yet! Type " + p + "ohayou to get your first ration."
	}

	// item not found
	if !getItem(itm) {
		return "I don't have that in stock."
	}

	// item cannot be purchased
	if !ITEM.Purchase {
		return "That's not for sale."
	}

	if USER.Ohayous < ITEM.Price*amt {
		return "You can't afford that."
	}

	// user is already at the limit for that item
	if ITEM.Limit > 0 && USER.Items[itm] >= ITEM.Limit {
		return fmt.Sprintf("You can't purchase any more of that. You can only have"+
			" %d %s", ITEM.Limit, itm)
	}

	// this purchase (presumeably batch purchase) would push them over the limit
	if ITEM.Limit > 0 && USER.Items[itm]+amt > ITEM.Limit {
		return fmt.Sprintf("You can't purchase that much. You can only have"+
			" %d %s", ITEM.Limit, itm)
	}

	// if pin is 4 digits
	if USER.Pin > 999 {
		go needPin(USER.Username, itm, channel, USER.Pin, amt)
		say(USER.Username, "You must enter your PIN to verify your purchase. "+
			"You have 10 seconds to do so or this purchase will not go through."+
			" NOTE: There is a bug right now. You may need to enter your pin here"+
			" more than once for it work.")
		return USER.Username + ": check your PM to verify your purchase."
	}

	USER.saveItem(itm, amt)

	if amt > 1 {
		return fmt.Sprintf("You purchased %d %ss for %d ohayous. "+
			"You have %d ohayous left.",
			amt, itm, ITEM.Price*amt, USER.Ohayous-(ITEM.Price*amt))
	} else {
		return fmt.Sprintf("You purchased %d %s for %d ohayous. "+
			"You have %d ohayous left.",
			amt, itm, ITEM.Price*amt, USER.Ohayous-(ITEM.Price*amt))
	}
}

func needPin(u, itm, channel string, upn, amt int) {
	for {
		select {
		case pn := <-getPin:
			if upn != pn.Pin && u == pn.Username {
				say(u, "Invalid PIN. Try again.")
			}
			if upn == pn.Pin && u == pn.Username {
				getUser(u)
				USER.saveItem(itm, amt)
				if amt > 1 {
					say(channel, fmt.Sprintf("%s purchased %d %ss for %d "+
						"ohayous. You have %d ohayous left.",
						USER.Username, amt, itm, ITEM.Price*amt,
						USER.Ohayous-(ITEM.Price*amt)))
				} else {
					say(channel, fmt.Sprintf("%s purchased %d %s for %d "+
						"ohayous. You have %d ohayous left.",
						USER.Username, amt, itm, ITEM.Price*amt,
						USER.Ohayous-(ITEM.Price*amt)))
				}
				return
			}
		case <-time.After(time.Second * 10):
			say(channel, fmt.Sprintf("%s failed to verify their purchase. This "+
				"has been reported as potential fraud.", u))
			return
		}
	}
}
