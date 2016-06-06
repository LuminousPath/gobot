package ohayou

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_quarry(m common.EmitMsg) {
	user, ok := GetUser(strings.ToLower(m.Nick))
	// user has never ohayou'd
	if !ok {
		say(m.Nick, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	// if nick is registered but not identified
	if user.Registered && !identified[user.Username] {
		say(user.Username, user.Username+": You must be identified with me to do "+
			"that. Make sure you are identified with the network and then type "+p+
			"identify.")
		return
	}

	// user does not have a quarry
	if user.Quarry.Installed == 0 {
		say(m.Nick, "You don't have any quarries yet. Keep saving!")
		return
	}

	inv := fmt.Sprintf("You have %d quarries and have mined these metals: ",
		user.Quarry.Installed)

	n := map[int][]string{}
	var a []int
	for k, v := range user.Quarry.Metals {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, amt := range a {
		for _, itm := range n[amt] {
			if amt == 0 {
				continue
			} else {
				inv += fmt.Sprintf("%d %s, ", amt, itm)
			}
		}
	}

	// PM the inv result with the trailing ", " cut off
	say(m.Nick, inv[:len(inv)-2])
}
