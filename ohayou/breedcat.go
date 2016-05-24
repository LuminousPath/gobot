package ohayou

import (
	"time"
)

func (u *User) breedCat() {
	// 2 - 4 hours random (7200, 14400)
	aloneTime := time.Duration(randNum(20, 40))
	breedTimer := time.NewTimer(aloneTime * time.Second)
	go func() {
		<-breedTimer.C
		if randNum(0, 10) <= 3 {
			say(u.Username, "Darn! Looks like your cats didn't mate, "+u.Username+
				". Try again tomorrow!")
		} else {
			say(u.Username, "Congratulations "+u.Username+"! Your cats "+
				"successfully mated! Amazingly it was born instantaneously! "+
				"You receive one cat.")
			u.SaveCat()
		}
		breedTimer.Stop()
	}()
}
