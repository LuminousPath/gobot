package ohayou

import (
	"fmt"
	"time"
)

func (u *User) PumpOil() {
	pumpTimer := time.NewTimer(6 * time.Second)
	go u.SetStatus("pumping", true)

	<-pumpTimer.C
	amt := randNum(1, 9) * u.Items["oilwell"]
	if u.Items["oilwell"] > 1 {
		say(u.Username, fmt.Sprintf("Your oil wells pumped %d barrels of crude oil.",
			amt))
	} else {
		say(u.Username, fmt.Sprintf("Your oil well pumped %d barrels of crude oil.",
			amt))
	}
	pumpTimer.Stop()
	u.SetStatus("pumping", false)
	u.SaveOil(amt)
}
