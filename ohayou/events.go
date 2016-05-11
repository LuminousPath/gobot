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
	catRand = time.Duration(randNum(2, 8))
	catTimer = time.NewTimer(catRand * time.Hour)
	go func() {
		<-catTimer.C
		for _, c := range chans {
			say(c, "A random cat appears! "+p+"use burger or pancake to adopt it!")
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
					say(c, user+" adopted the cat!")
				}
				USER.saveItem("cat", 1)
			}
			canAdoptCat = false
			return
		case <-time.After(time.Second * 10):
			for _, c := range chans {
				say(c, "The cat wanders off and disappears...")
			}
			canAdoptCat = false
			return
		}
	}
}
