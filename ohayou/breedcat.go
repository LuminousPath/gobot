package ohayou

import (
	"fmt"
	"time"
)

func (u *User) breedCat() {
	// 2 - 4 hours random (7200, 14400)
	aloneTime := time.Duration(randNum(7200, 14400))
	breedTimer := time.NewTimer(aloneTime * time.Second)

	go u.SetStatus("breeding", true)

	<-breedTimer.C

	if randNum(0, 10) <= 3 {
		say(u.Username, "Darn! Looks like your cats didn't mate, "+u.Username+
			". Maybe next time!")
	} else {
		litter := randNum(2, 7) * u.Items["cattery"]
		say(u.Username, fmt.Sprintf("Congratulations "+u.Username+"! Your cats "+
			"successfully mated! Amazingly, the cats were born instantaneously! "+
			"You receive %d cats.", litter))
		u.SaveCat(litter)
	}
	breedTimer.Stop()
	u.SetStatus("breeding", false)
}
