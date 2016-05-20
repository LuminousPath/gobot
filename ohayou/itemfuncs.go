package ohayou

import (
	"time"
)

var (
	itemFuncs = map[string]func(*User, string) string{
		"saveBottle":  saveBottle,
		"dragonDildo": dragonDildo,
		"adoptCat":    adoptCat,
		"fortune":     fortune}
)

func saveBottle(u *User, itm string) string {
	item, _ := GetItem("bottle")
	go u.SaveItem(item, 1)
	return ""
}

func dragonDildo(u *User, itm string) string {
	if u.LastUsed[itm].In(est).Format("20060102") >= time.Now().In(est).Format("20060102") {
		return " but is all out of lube."
	}
	go u.SetLastUsed(itm)
	go u.ResetLast()
	return " and shortly thereafter feels good enough to " + p + "ohayou again."
}

func adoptCat(u *User, itm string) string {
	if canAdoptCat {
		catAdopt <- u.Username
	}
	return ""
}

func fortune(u *User, itm string) string {
	t := time.Now()
	if u.LastUsed[itm].In(est).Format("20060102") >= t.In(est).Format("20060102") {
		return "- here's today's fortune again: " + u.Fortune
	}
	go u.SetLastUsed(itm)
	// see fortune.go -- GetFortune() saves and return fortune
	return "- " + u.GetFortune()
}