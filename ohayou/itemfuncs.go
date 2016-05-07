package ohayou

import (
	"log"
	"time"
)

var itemFuncs = map[string]func(*User, string) string{
	"saveBottle":  saveBottle,
	"dragonDildo": dragonDildo}

func saveBottle(user *User, itemName string) string {
	saveItem(user, "bottle", 1)
	return ""
}

func dragonDildo(user *User, itemName string) string {
	t := time.Now()
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("err: ", err.Error())
	}

	if user.LastUsed[itemName].Format("20060102") >= t.In(est).Format("20060102") {
		return " but is all out of lube"
	}

	go setLastUsed(user, itemName)
	go resetLast(user)

	return " and shortly thereafter feels good enough to " + p + "ohayou again."
}
