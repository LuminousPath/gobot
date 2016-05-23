package ohayou

import (
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_stats(m common.EmitMsg) {
	var to string
	if isPM(m) {
		to = m.Nick
	} else {
		to = m.Channel
	}

	user, ok := GetUser(strings.ToLower(m.Nick))
	if !ok {
		say(to, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	go user.Stats()
}
