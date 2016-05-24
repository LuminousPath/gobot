package ohayou

import "time"

type User struct {
	Username       string
	Last           time.Time
	Ohayous        int
	CumOhayous     int `bson:"cumOhayous"`
	StealSuccess   int `bson:"stealSuccess"`
	StealFail      int `bson:"stealFail"`
	StolenFrom     int `bson:"stolenFrom"`
	StolenOhayous  int `bson:"stolenOhayous"`
	OhayousStolen  int `bson:"ohayousStolen"`
	Probation      time.Time
	ProbationCount int `bson:"probationCount"`
	TimesOhayoued  int `bson:"timesOhayoued"`
	Items          map[string]int
	ItemMultiply   map[string]int `bson:"itemMultiply"`
	Equipped       map[string]Item
	LastUsed       map[string]time.Time `bson:"lastUsed"`
	Registered     bool
	Fortune        string
	Vault
}

type Item struct {
	Name          string
	Desc          string
	Price         int
	Add           int
	Multiply      int
	Multiplies    string
	Defense       int
	Limit         int
	Acrelimit     int
	Useable       bool
	Consume       bool
	Effect        string
	HasFunction   string `bson:"hasFunction"`
	Purchase      bool
	Category      string
	EquipCategory string
}

type Vault struct {
	Installed bool
	Level     int
	Ohayous   int
	Last      time.Time
}

type TimeZone struct {
	l *time.Location
	e error
}

type UserOhayous struct {
	Username string
	Ohayous  int
}
