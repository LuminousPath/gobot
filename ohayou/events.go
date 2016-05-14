package ohayou

import (
	"time"
)

var (
	// cat event stuff
	catTimer      *time.Timer
	catRand       time.Duration
	catAdoptTimer *time.Timer
	canAdoptCat   bool
	catAdopt      = make(chan string)

	eventFuncs = map[string]func(){
		"catEvent": catEvent}
)

func startEvents() {
	for e := range eventFuncs {
		startEvent := eventFuncs[e]
		go startEvent()
	}
}

func catEvent() {
	// 2 - 8 hours random
	catRand = time.Duration(randNum(7200, 28800))
	catTimer = time.NewTimer(catRand * time.Second)
	go func() {
		<-catTimer.C
		for _, c := range chans {
			say(c, "A random cat appears! "+p+"use burger or "+p+"use pancake to "+
				"adopt it!")
		}
		catTimer.Stop()
		go waitCatAdopt()
		catEvent()
	}()
}

func waitCatAdopt() {
	canAdoptCat = true
	for {
		select {
		case user := <-catAdopt:
			if getUser(user) {
				for _, c := range chans {
					say(c, "And adopts the cat!")
				}
				USER.saveItem("cat", 1)
			}
			canAdoptCat = false
			return
		case <-time.After(time.Second * 15):
			for _, c := range chans {
				say(c, "The cat wanders off and disappears...")
			}
			canAdoptCat = false
			return
		}
	}
}
