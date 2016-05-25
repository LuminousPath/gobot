package ohayou

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_inventory(m common.EmitMsg) {
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

	// user does not have any item
	if len(user.Items) == 0 {
		say(m.Nick, "You don't have any items yet. Keep saving!")
		return
	}

	inv := fmt.Sprintf("You have: %d ohayous, ", user.Ohayous)
	// if user has a vault, make it the first item and show ohayous in it
	if user.Vault.Installed {
		inv += fmt.Sprintf("a Level %d vault (%d/%d ohayous), ",
			user.Vault.Level+1, user.Vault.Ohayous,
			int(math.Pow(10, 3+float64(user.Vault.Level))))
	}

	n := map[int][]string{}
	var a []int
	for k, v := range user.Items {
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
			} else if itm == "vault" {
				continue
			} else if amt > 1 {
				inv += fmt.Sprintf("%d %ss, ", amt, itm)
			} else {
				inv += fmt.Sprintf("%d %s, ", amt, itm)
			}
		}
	}

	// PM the inv result with the trailing ", " cut off
	say(m.Nick, inv[:len(inv)-2])
}
