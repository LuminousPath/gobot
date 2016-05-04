package ohayou

import "time"

type User struct {
	Username      string
	Last          time.Time
	Ohayous       int
	CumOhayous    int
	Add           int
	TimesOhayoued int
	Items         map[string]int
	ItemMultiply  map[string]int
}

type Item struct {
	Name       string
	Desc       string
	Price      int
	Add        int
	Multiply   int
	Multiplies bool
	Defense    int
	Limit      int
	Useable    bool
	Consume    bool
	Effect     string
	Purchase   bool
	Category   string
}
