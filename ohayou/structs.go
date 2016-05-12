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
	ItemMultiply  map[string]int       `bson:"itemMultiply"`
	LastUsed      map[string]time.Time `bson:"lastUsed"`
	Pin           int
}

type Item struct {
	Name        string
	Desc        string
	Price       int
	Add         int
	Multiply    int
	Multiplies  string
	Defense     int
	Limit       int
	Useable     bool
	Consume     bool
	Effect      string
	HasFunction string `bson:"hasFunction"`
	Purchase    bool
	Category    string
}

type TimeZone struct {
	l *time.Location
	e error
}

type UserOhayous struct {
	Username string
	Ohayous  int
}

type SubmitPin struct {
	Username string
	Pin      int
}
