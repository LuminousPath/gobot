package catfact

import (
	"encoding/json"
	"github.com/mferrera/go-ircevent"
	"log"
	"net/http"
)

type CatFact struct {
	Fact    []string `json:"facts"`
	Success string   `json:"success"`
}

func fact() string {
	api := "http://catfacts-api.appspot.com/api/facts?number=1"

	resp, err := http.Get(api)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	cat := CatFact{}

	err = decoder.Decode(&cat)
	if err != nil {
		log.Println(err)
	}

	if cat.Success == "true" {
		return cat.Fact[0]
	} else {
		return "There was a problem."
	}
}

func Run(b **irc.Connection, p, cmd, channel string, word []string, admin bool) bool {
	bot := *b
	say := bot.Privmsg
	var fact string = fact()

	if word[0] == p+"cat" || word[0] == p+"catfact" {
		say(channel, fact)
	}

	return true
}
