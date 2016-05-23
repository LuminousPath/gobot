package ohayou

import (
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_ohayou(m common.EmitMsg) {
	if isPM(m) {
		say(m.Nick, "You can only do that in a channel I'm in.")
		return
	}

	say(m.Channel, Ohayou(strings.ToLower(m.Nick)))
}
