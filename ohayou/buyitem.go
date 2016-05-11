package ohayou

import (
	"fmt"
)

func buyItem(nick, itm string, amt int) string {
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
