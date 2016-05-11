package ohayou

import (
	"time"
)

var itemFuncs = map[string]func(*User, string) string{
	"saveBottle":  saveBottle,
	"dragonDildo": dragonDildo,
	"adoptCat":    adoptCat}

func saveBottle(u *User, itm string) string {
	go u.saveItem("bottle", 1)
	return ""
}

func dragonDildo(u *User, itm string) string {
	t = time.Now()
	if u.LastUsed[itm].Format("20060102") >= t.In(est).Format("20060102") {
		return " but is all out of lube"
	}
	go u.setLastUsed(itm)
	go u.resetLast()
	return " and shortly thereafter feels good enough to " + p + "ohayou again."
}

func adoptCat(u *User, itm string) string {
	if canAdoptCat {
		catAdopt <- u.Username
	}
	return ""
}
