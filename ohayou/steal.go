package ohayou

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// percent chance of success
	stealOhayouSuccess = 36
	stealCatSuccess    = 24

	stealFineMin   int     = 5
	stealFinePct   float64 = 0.16
	stealAmountPct float64 = 0.07
)

func (t *User) stealFine() int {
	return stealFineMin + int(float64(t.Ohayous)*stealFinePct)
}

func (v *User) stealAmount() int {
	return int(float64(v.Ohayous) * stealAmountPct)
}

// t = thief, v = victim
func (t *User) StealFrom(v User, channel, nickRaw, vicRaw string) {
	// formatted time to the minute
	if t.Probation.In(est).Format("200601021504") >=
		time.Now().In(est).Format("200601021504") {

		say(channel, fmt.Sprintf("%s: you are still on probation from your last "+
			"theft attempt. Your probation expires on %v EST.",
			nickRaw, t.Probation.In(est).Format("Jan 2 15:04")))
		return
	}
	stealOhayouChance := randNum(0, 100)
	defense := 0
	var stealCatChance int
	if v.Items["cat"] == 0 {
		stealCatChance = 0
	} else {
		stealCatChance = randNum(0, 100)
	}

	if v.Equipped != nil {
		for _, name := range v.Equipped {
			defense += name.Defense
		}
		stealOhayouChance -= int(defense / 9)
		stealCatChance -= int(defense / 7)
	}

	if t.TimesOhayoued < 5 {
		say(channel, nickRaw+": you haven't ohayou'd enough to do that yet!")
		return
	}

	if t.Ohayous < stealFineMin {
		say(channel, fmt.Sprintf("%s: you don't have enough ohayous to steal. "+
			"You need at least %d.", nickRaw, stealFineMin))
		return
	}

	if v.Ohayous == 0 {
		say(channel, fmt.Sprintf("%s attempts to steal from %s but %s "+
			"doesn't have any ohayous! %s is fined %d ohayous and placed "+
			"on probation for 24 hours.",
			nickRaw, vicRaw, vicRaw, nickRaw, t.stealFine()))
		t.SaveFailSteal(t.stealFine())
		return
	}

	// if steal fails on ohayou and cat
	if stealOhayouChance > stealOhayouSuccess && stealCatChance > stealCatSuccess {
		say(channel, fmt.Sprintf("%s attempts to steal from %s but is caught! "+
			"%s is fined %d ohayous and is placed on probation for 24 hours.",
			nickRaw, vicRaw, nickRaw, t.stealFine()))
		t.SaveFailSteal(t.stealFine())
		return
	}

	// if steal fails on cat, succeeds on ohayous
	if stealOhayouChance <= stealOhayouSuccess && stealCatChance > stealCatSuccess {
		say(channel, fmt.Sprintf("%s attempts to steal from %s and succeeds! "+
			"%s steals %d ohayous from %s.",
			nickRaw, vicRaw, nickRaw, v.stealAmount(), vicRaw))
		SaveSuccessSteal(t, v, 0, v.stealAmount())
		return
	}

	// if steal fails on ohayous, succeeds on cat
	if stealOhayouChance > stealOhayouSuccess && stealCatChance <= stealCatSuccess {
		say(channel, fmt.Sprintf("%s attempts to steal from %s and succeeds! "+
			"%s steals a cat from %s.",
			nickRaw, vicRaw, nickRaw, vicRaw))
		SaveSuccessSteal(t, v, 1, 0)
		return
	}

	// if both ohayou and cat thefts succeed
	if stealOhayouChance <= stealOhayouSuccess && stealCatChance <= stealCatSuccess {
		say(channel, fmt.Sprintf("%s attempts to steal from %s and succeeds! "+
			"%s steals a cat and %d ohayous from %s.",
			nickRaw, vicRaw, nickRaw, v.stealAmount(), vicRaw))
		SaveSuccessSteal(t, v, 1, v.stealAmount())
		return
	}
}

func SaveSuccessSteal(t *User, v User, cat, ohy int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	saveThief := bson.M{"$inc": bson.M{
		"ohayous":       ohy,
		"items.cat":     cat,
		"stealSuccess":  1,
		"stolenOhayous": ohy}}
	saveVictim := bson.M{"$inc": bson.M{
		"ohayous":       -ohy,
		"items.cat":     -cat,
		"stolenFrom":    1,
		"ohayousStolen": ohy}}

	err := q.Update(bson.M{"username": t.Username}, saveThief)
	if err != nil {
		log.Println("SaveItem: " + err.Error())
	}

	err = q.Update(bson.M{"username": v.Username}, saveVictim)
	if err != nil {
		log.Println("SaveItem: " + err.Error())
	}
}

func (t *User) SaveFailSteal(fine int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{
		"ohayous":        t.Ohayous - fine,
		"probation":      time.Now().Add(time.Duration(24) * time.Hour).In(est),
		"probationCount": t.ProbationCount + 1,
		"stealFail":      t.StealFail + 1}}

	err := q.Update(bson.M{"username": t.Username}, save)
	if err != nil {
		log.Println("SaveItem: " + err.Error())
	}
}
