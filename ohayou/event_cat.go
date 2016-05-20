package ohayou

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	catTimer      *time.Timer
	catRand       time.Duration
	catAdoptTimer *time.Timer
	canAdoptCat   bool
	catAdopt      = make(chan string)
)

func catEvent() {
	// 2 - 8 hours random (7200, 28800)
	catRand = time.Duration(randNum(7200, 28800))
	catTimer = time.NewTimer(catRand * time.Second)
	go func() {
		<-catTimer.C
		for _, c := range chans {
			say(c, "A stray cat appears! "+p+"use burger or "+p+"use pancake to "+
				"adopt it!")
		}
		catTimer.Stop()
		go waitCatAdopt()
		catEvent()
	}()
}

func waitCatAdopt() {
	canAdoptCat = true
	select {
	case r := <-catAdopt:
		user, ok := GetUser(r)
		if ok {
			for _, c := range chans {
				say(c, user.Username+" adopts the cat!")
			}
			user.SaveCat()
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

func (u *User) SaveCat() {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$inc": bson.M{
		"items.cat": 1}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveItem: " + err.Error())
	}
}
