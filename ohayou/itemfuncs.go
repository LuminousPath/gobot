package ohayou

import (
	"time"
)

var (
	itemFuncs = map[string]func(*User, string) string{
		"saveBottle":      saveBottle,
		"dragonDildo":     dragonDildo,
		"adoptCat":        adoptCat,
		"fortune":         fortune,
		"makeVault":       makeVault,
		"attemptBreedCat": attemptBreedCat,
		"startMining":     startMining,
		"startPumping":    startPumping}
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
	if u.LastUsed[itm].In(est).Format("20060102") >= time.Now().In(est).Format("20060102") {
		return "- here's today's fortune again: " + u.Fortune
	}
	go u.SetLastUsed(itm)
	// see fortune.go -- GetFortune() saves and return fortune
	return "- " + u.GetFortune()
}

func makeVault(u *User, itm string) string {
	go u.InstallVault()
	return ""
}

func attemptBreedCat(u *User, itm string) string {
	if u.Items["cat"] < 2 {
		return " but doesn't have any cats to breed! What are you doing! You need " +
			"at least two cats to do that."
	}

	if u.Status["breeding"] {
		return " but already has cats in there! You must wait until they are finished."
	}

	go u.breedCat()
	return " for a few hours."
}

func startMining(u *User, itm string) string {
	if u.Status["mining"] {
		return " but is already mining! Wait until it's finished and try again."
	}

	go u.Mine()
	return " for a few hours."
}

func startPumping(u *User, itm string) string {
	if u.Status["pumping"] {
		return " but is already pumping oil! Wait until it's finished and try again."
	}

	go u.PumpOil()
	return " for a few hours."
}
