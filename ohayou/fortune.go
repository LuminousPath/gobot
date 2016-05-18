package ohayou

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

type FortuneHolder struct {
	Fortune [500]string
}

var allFortunes FortuneHolder // stores all fortunes

func (u *User) GetFortune() string {
	newFortune := allFortunes.Fortune[randNum(0, 500)]
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{
		"fortune": newFortune}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("GetFortune: " + err.Error())
	}
	// returns to be printed in channel
	return newFortune
}

func fillFortunes() {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C("extra")

	err := q.FindId(bson.ObjectIdHex("573776916ef0ad54ad2d08a2")).One(&allFortunes)
	if err != nil {
		log.Println("fillFortunes: " + err.Error())
	}
}
