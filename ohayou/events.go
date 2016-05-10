package ohayou

import (
	"log"
	"time"

	"github.com/mferrera/go-ircevent"
)

// keeps a list of events herein to check if they are active
var eventFuncs = map[string]func(*irc.Connection){
	"catEvent": catEvent}

func startEvents(bot *irc.Connection) {
	for e := range eventFuncs {
		startEvent := eventFuncs[e]
		go startEvent(bot)
	}
}

func catEvent(bot *irc.Connection) {
	quitChan := make(chan bool)

	howLong := time.Duration(randNum(2, 8))
	// stop timer after event has fired
	go func() {
		<-time.After(howLong * time.Second)
		close(quitChan)
	}()

	// fire once per second
	t := time.NewTicker(howLong * time.Second)
	for {
		select {
		case <-t.C:
		case <-quitChan:
			t.Stop()
			return
		}
		// restart event
		catEvent(bot)
	}
}
