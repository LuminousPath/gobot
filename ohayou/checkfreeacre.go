package ohayou

import "math"

func (u *User) FreeAcre(amt int) bool {
	for itm, uAmt := range u.Items {
		item, _ := GetItem(itm)
		if item.Acrelimit <= 0 {
			continue
		}

		usedAcres := math.Ceil(float64(uAmt) / float64(item.Acrelimit))
		if (u.Items["acre"] - int(usedAcres)) < amt {
			return false
		}
	}
	return true
}
