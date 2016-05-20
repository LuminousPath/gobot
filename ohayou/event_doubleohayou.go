package ohayou

import (
	"time"
)

var (
	doubleOhayou       bool
	doubleOhayouRand   time.Duration
	doubleOhayouTimer  *time.Timer
	doubleOhayouFRand  time.Duration
	doubleOhayouFTimer *time.Timer
)

func doubleOhayouEvent() {
	// 8 - 36 hours random (28800, 129600)
	doubleOhayouRand = time.Duration(randNum(43200, 129600))
	doubleOhayouTimer = time.NewTimer(doubleOhayouRand * time.Second)
	go func() {
		<-doubleOhayouTimer.C
		for _, c := range chans {
			say(c, "ERROR: Ohayou distributor is malfunctioning.")
		}
		doubleOhayouTimer.Stop()
		doubleOhayou = true
		doubleOhayouFRand = time.Duration(randNum(2, 10))
		doubleOhayouFTimer = time.NewTimer(doubleOhayouFRand * time.Minute)
		go func() {
			<-doubleOhayouFTimer.C
			doubleOhayou = false
			for _, c := range chans {
				say(c, "Technicians have fixed the ohayou distributor. "+
					"It should be working as normal now.")
			}
			doubleOhayouFTimer.Stop()
		}()
		doubleOhayouEvent()
	}()
}
