package ohayou

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

func (u *User) Equip(itm string) string {
	// if nick is registered but not identified
	if u.Registered && !identified[u.Username] {
		return u.Username + ": You must be identified with me to do that. Make sure " +
			"you are identified with the network and then type " + p + "identify."
	}

	item, ok := GetItem(itm)

	if !ok {
		return u.Username + ": That isn't an item."
	}

	if u.Items[itm] <= 0 {
		return u.Username + ": You don't have that item."
	}

	if item.EquipCategory == "" {
		return u.Username + ": That item can't be equipped."
	}
	// if item is already equipped
	if u.Equipped[item.EquipCategory].Name == item.Name {
		return u.Username + ": You already have that item equipped."
	}

	// if nothing is equipped there
	if u.Equipped[item.EquipCategory].Name == "" {
		u.SaveEquip(item)
		return u.Username + " equipped " + itm + "."
	}

	// if something else is already equipped in the item's slot
	if u.Equipped[item.EquipCategory].Name != item.Name {
		u.SaveEquip(item)
		return u.Username + " unequipped " + u.Equipped[item.EquipCategory].Name +
			" from " + item.EquipCategory + " and equipped " + item.Name + "."
	}

	return ""
}

func (u *User) Unequip(itm string) string {
	// if nick is registered but not identified
	if u.Registered && !identified[u.Username] {
		return u.Username + ": You must be identified with me to do that. Make sure " +
			"you are identified with the network and then type " + p + "identify."
	}

	item, ok := GetItem(itm)

	if !ok {
		return u.Username + ": That isn't an item."
	}

	// if item is already equipped
	if u.Equipped[item.EquipCategory].Name == item.Name {
		u.SaveUnequip(item.EquipCategory)
		return u.Username + " unequipped " + item.Name + " from " +
			item.EquipCategory + "."
	}

	// if nothing is equipped there
	if u.Equipped[item.EquipCategory].Name == "" {
		return u.Username + ": That item isn't equipped."
	}

	return ""
}

func (u *User) SaveEquip(item Item) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{"equipped." + item.EquipCategory: item}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveOhayous: " + err.Error())
	}
}

func (u *User) SaveUnequip(equipcategory string) {
	s := session.Copy()
	defer s.Close()
	q := s.DB(dbName).C(ohyCol)

	save := bson.M{"$set": bson.M{"equipped." + equipcategory: nil}}

	err := q.Update(bson.M{"username": u.Username}, save)
	if err != nil {
		log.Println("saveOhayous: " + err.Error())
	}
}
