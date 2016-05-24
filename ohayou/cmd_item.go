package ohayou

import (
	"fmt"
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_item(m common.EmitMsg) {
	var to string
	if isPM(m) {
		to = m.Nick
	} else {
		to = m.Channel
	}

	// if no argument is given for the command
	if !hasArgs(m.Word) {
		say(to, "Gives information about a specific item. Usage: "+p+"item <itemname>")
		return
	}

	item, ok := GetItem(strings.ToLower(m.Word[1]))
	// if the queried item doesn't
	if !ok {
		say(to, "I don't carry that item.")
		return
	}

	// if the item cannot be purchased
	if !item.Purchase {
		say(to, fmt.Sprintf("%s: %s. Cannot be purchased.", item.Name, item.Desc))
		return
	}

	// if theres a acreage limit, list it
	if item.Acrelimit > 0 {
		say(to, fmt.Sprintf("%s: %s - Price: %d ohayous. Limited to %d per acre.",
			item.Name, item.Desc, item.Price, item.Acrelimit))
		return
	}

	// say information about the item
	say(to, fmt.Sprintf("%s: %s - Price: %d ohayous.", item.Name, item.Desc, item.Price))
}
