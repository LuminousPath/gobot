package ohayou

import "time"

type User struct {
	Username      string
	Last          time.Time
	Ohayous       int
	CumOhayous    int `bson:"cumOhayous"`
	Add           int
	TimesOhayoued int `bson:"timesOhayoued"`
	Items         map[string]int
	ItemMultiply  map[string]int `bson:"itemMultiply"`
}

type Item struct {
	Name       string
	Desc       string
	Price      int
	Add        int
	Multiply   int
	Multiplies string
	Defense    int
	Limit      int
	Useable    bool
	Consume    bool
	Effect     string
	Purchase   bool
	Category   string
}
