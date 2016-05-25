package ohayou

import (
	"fmt"
	"log"
	"math"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func (u *User) Deposit(amt int) string {
	// if nick is registered but not identified
	if u.Registered && !identified[u.Username] {
		return u.Username + ": You must be identified with me to do that. Make sure " +
			"you are identified with the network and then type " + p + "identify."
	}

	cap := int(math.Pow(10, 3+float64(u.Vault.Level)))
	if !u.Vault.Installed {
		return u.Username + ": You don't have a vault yet."
	}

	if u.Vault.Ohayous >= cap {
		return u.Username + ": Your vault is full. Consider upgrading it."
	}
	if u.Vault.Last.In(est).Format("20060102") >= time.Now().In(est).Format("20060102") {
		return u.Username + ": You've already opened your vault once today. Due to " +
			"security concerns you cannot open it again."
	}

	if u.Ohayous < amt {
		return u.Username + ": You don't have that many ohayous."
	}

	if (u.Vault.Ohayous + amt) > (cap - u.Vault.Ohayous) {
		return u.Username + ": That's more than your vault can hold. Double-check " +
			"your numbers or purchase an upgrade."
	}

	go u.depositOhayous(amt)
	return fmt.Sprintf("%s deposited %d ohayous to their vault.", u.Username, amt)
}

func (u *User) depositOhayous(amt int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{
		"ohayous":       u.Ohayous - amt,
		"vault.ohayous": u.Vault.Ohayous + amt,
		"vault.last":    time.Now().In(est)}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("SaveItem: " + err.Error())
	}
}

func (u *User) Withdraw(amt int) string {
	// if nick is registered but not identified
	if u.Registered && !identified[u.Username] {
		return u.Username + ": You must be identified with me to do that. Make sure " +
			"you are identified with the network and then type " + p + "identify."
	}

	if !u.Vault.Installed {
		return u.Username + ": You don't have a vault yet."
	}

	if u.Vault.Ohayous == 0 {
		return u.Username + ": You don't have any ohayous in your vault."
	}

	if u.Vault.Last.In(est).Format("20060102") >= time.Now().In(est).Format("20060102") {
		return u.Username + ": You've already opened your vault once today. " +
			"According to vault security protocol you cannot open it again until" +
			" tomorrow."
	}

	if (u.Vault.Ohayous - amt) < 0 {
		return u.Username + ": You don't have that many ohayous in your vault."
	}

	go u.withdrawOhayous(amt)
	return fmt.Sprintf("%s withdrew %d ohayous from their vault.", u.Username, amt)
}

func (u *User) withdrawOhayous(amt int) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{
		"ohayous":       u.Ohayous + amt,
		"vault.ohayous": u.Vault.Ohayous - amt,
		"vault.last":    time.Now().In(est)}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("SaveItem: " + err.Error())
	}
}
