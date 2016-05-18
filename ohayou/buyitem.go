package ohayou

import (
	"fmt"
	"time"
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

	if u.Ohayous < item.Price*amt {
		return "You can't afford that."
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

	// if pin is 4 digits
	if u.Pin > 999 {
		say(u.Username, "You must enter your PIN to verify your purchase. "+
			"You have 10 seconds to do so or this purchase will not go through."+
			" NOTE: There is a bug right now. You may need to enter your pin here"+
			" more than once for it work.")
		go u.EnterPin(channel, item, amt)
		return u.Username + ": check your PM to verify your purchase."
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

func (u *User) EnterPin(channel string, item Item, amt int) {
	select {
	case r := <-sendPin:
		if u.Pin != r.Pin && u.Username == r.Username {
			say(u.Username, "Invalid PIN. Try again.")
		}
		if u.Pin == r.Pin && u.Username == r.Username {
			u.SaveItem(item, amt)
			if amt > 1 {
				say(channel, fmt.Sprintf("%s purchased %d %ss for %d "+
					"ohayous. You have %d ohayous left.",
					u.Username, amt, item.Name, item.Price*amt,
					u.Ohayous-(item.Price*amt)))
			} else {
				say(channel, fmt.Sprintf("%s purchased %d %s for %d "+
					"ohayous. You have %d ohayous left.",
					u.Username, amt, item.Name, item.Price*amt,
					u.Ohayous-(item.Price*amt)))
			}
			return
		}
	case <-time.After(time.Second * 10):
		say(channel, fmt.Sprintf("%s failed to verify their purchase. This "+
			"has been reported as potential fraud.", u))
		return
	}
}
