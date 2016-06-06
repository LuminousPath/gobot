package ohayou

import (
	"fmt"
	"time"
)

// %chance to receive these items from mining
var metals map[string]int = map[string]int{
	"aluminum": 100,
	"iron":     95,
	"titanium": 85,
	"copper":   60,
	"lead":     50,
	"tin":      40,
	"uranium":  30,
	"silver":   25,
	"platinum": 20,
	"gold":     15,
}

func (u *User) Mine() {
	mineTimer := time.NewTimer(8 * time.Hour)
	yield := make(map[string]int)
	sayYield := "You mined "
	go u.SetStatus("mining", true)

	<-mineTimer.C
	for mtl, ch := range metals {
		chance := randNum(0, 100)
		amt := (1 + metals[mtl]/10) * u.Quarry.Installed
		if chance < ch {
			yield[mtl] = amt
			sayYield += fmt.Sprintf("%d %s, ", amt, mtl)
		}
	}
	say(u.Username, sayYield[:len(sayYield)-2])
	mineTimer.Stop()
	u.SetStatus("mining", false)
	u.SaveMetals(yield)
}
