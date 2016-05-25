package ohayou

import (
	"strings"
	"time"

	"github.com/mferrera/go-ircevent"
)

var policeProtected = make(map[string][]int)

func (u *User) StationPolice(b *irc.Connection) {
	if _, protected := policeProtected[u.Username]; protected {
		return
	}

	say(u.Username, "Ohayou Police here. Looks like you were just the victim of a "+
		"robbery. If you report it, we can station one of our officers nearby for a "+
		"couple of hours. It'll reduce the chance of it happening again. Type "+p+
		"report if you're interested.")

	// 1 minute to respond
	responseTimer := time.NewTimer(60 * time.Second)
	// take the ID of the callback to later remove it
	var responseId int
	responseId = b.AddCallback("PRIVMSG", func(e *irc.Event) {
		word := strings.Split(e.Message(), " ")
		if strings.ToLower(e.Nick) == u.Username && word[0] == p+"report" {
			responseTimer.Stop()
			say(u.Username, "Alright, we'll watch over you for a few hours.")
			go u.SetPoliceProtection()
			b.RemoveCallback("PRIVMSG", responseId)
			return
		}
	})

	// if no response
	go func() {
		<-responseTimer.C
		say(u.Username, "Guess you're not interested. Good luck out there.")
		b.RemoveCallback("PRIVMSG", responseId)
		responseTimer.Stop()
	}()
}

func (u *User) SetPoliceProtection() {
	// substract defensive items from the base chance of being stolen from
	userStealChance := stealOhayouSuccess - int(u.UserDefense()/9)
	userStealCatChance := stealCatSuccess - int(u.UserDefense()/7)

	// grab 90% of this -- police protection begins at leaving only a 10% chance of being
	// stolen from and slowly decreases
	modOhayouChance := int(float64(userStealChance) * 0.9)
	modCatChance := int(float64(userStealCatChance) * 0.9)

	// add it to the protected users hashmap
	policeProtected[u.Username] = append(policeProtected[u.Username], modOhayouChance)
	policeProtected[u.Username] = append(policeProtected[u.Username], modCatChance)

	// 90% of the user's chance of being stolen from after defensive items are factored in
	// divided by 4, since we remove a quarter of the police added defense every hour
	decOhayouChance := int(int(float64(userStealChance)*0.9) / 4)
	decCatChance := int(int(float64(userStealCatChance)*0.9) / 4)

	// set the ticker
	remProtection := time.NewTicker(1 * time.Hour)

	// remove some protection every hour until it's all gone
	for _ = range remProtection.C {
		policeProtected[u.Username][0] -= decOhayouChance
		policeProtected[u.Username][1] -= decCatChance
		// if all protection is gone, stop timer and let user know, delete from
		// map of protected persons
		if policeProtected[u.Username][0] == 0 {
			say(u.Username, "Ohayou Police here. We're leaving the vicinity now. "+
				"Good luck.")
			remProtection.Stop()
			delete(policeProtected, u.Username)
			break
		}
	}
}
