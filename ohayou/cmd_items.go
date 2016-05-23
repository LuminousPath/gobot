package ohayou

import (
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_items(m common.EmitMsg) {
	var to string
	if isPM(m) {
		to = m.Nick
	} else {
		to = m.Channel
	}

	if !hasArgs(m.Word) {
		say(to, "Type "+p+"items <category> to get a list of items by category. "+
			"Categories: "+strings.Join(append(itemCtgs), ", ")+".")
		return
	}

	itemsInCtg := ItemCategory(strings.ToLower(m.Word[1]))
	for _, itm := range itemsInCtg {
		say(m.Nick, itm)
	}
}
