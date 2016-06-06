package ohayou

import (
	"fmt"
)

func (u *User) Buy(channel, itm string, amt int) string {
	item, ok := GetItem(itm)
	// item not found
	if !ok {
		return "I don't have that in stock."
	}

	// item cannot be purchased
	if !item.Purchase {
		return "That's not for sale."
	}

	if amt <= 0 {
		return "That's not a valid quantity."
	}

	if u.Ohayous < item.Price*amt {
		return "You can't afford that."
	}

	if item.Acrelimit > 0 && (u.Items[itm]+amt) > (item.Acrelimit*u.Items["acre"]) {
		return fmt.Sprintf("You need more land to purchase more of that! You can "+
			"only have %d %ss per acre and you have %d %ss and %d acre(s).",
			item.Acrelimit, itm, u.Items[itm], itm, u.Items["acre"])
	}

	// user is already at the limit for that item
	if item.Limit > 0 && u.Items[itm] >= item.Limit {
		return fmt.Sprintf("You can't purchase any more of that. You can only have"+
			" %d %s", item.Limit, itm)
	}

	// this purchase (presumeably batch purchase) would push them over the limit
	if item.Limit > 0 && u.Items[itm]+amt > item.Limit {
		return fmt.Sprintf("You can't purchase that much. You can only have"+
			" %d %s", item.Limit, itm)
	}

	// if nick is registered but not identified
	if u.Registered && !identified[u.Username] {
		return u.Username + ": You must be identified with me to do that. Make sure " +
			"you are identified with the network and then type " + p + "identify."
	}

	if item.NeedsAcre && !u.FreeAcre(amt) {
		return u.Username + ": That item requires its own acre, and you do not have " +
			"an empty acre."
	}

	u.SaveItem(item, amt)

	if amt > 1 {
		return fmt.Sprintf("You purchased %d %ss for %d ohayous. "+
			"You have %d ohayous left.",
			amt, itm, item.Price*amt, u.Ohayous-(item.Price*amt))
	} else {
		return fmt.Sprintf("You purchased %d %s for %d ohayous. "+
			"You have %d ohayous left.",
			amt, itm, item.Price*amt, u.Ohayous-(item.Price*amt))
	}
}
