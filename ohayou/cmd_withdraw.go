package ohayou

import (
	"strconv"
	"strings"

	"github.com/mferrera/gobot/common"
)

func cmd_withdraw(m common.EmitMsg) {
	var to string
	if isPM(m) {
		to = m.Nick
	} else {
		to = m.Channel
	}

	// if no argument is given for the command
	if !hasArgs(m.Word) {
		say(to, "Withdraws ohayous from your vault. Usage: "+p+"withdraw <num> -- "+
			"withdraws <num> ohayous. Your vault can only be opened once per day "+
			"due to its security protocol.")
		return
	}

	user, ok := GetUser(strings.ToLower(m.Nick))
	if !ok {
		say(to, "You can't do that because you haven't ohayou'd yet! Type "+p+
			"ohayou to get your first ration.")
		return
	}

	amt, err := strconv.Atoi(m.Word[1])
	if err != nil {
		say(to, "You didn't give a valid quantity. Usage: "+p+"withdraw <num> will "+
			"withdraw <num> ohayous to your vault.")
		return
	}

	say(to, user.Withdraw(amt))
}
