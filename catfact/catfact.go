package catfact

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mferrera/go-ircevent"
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

	if resp.StatusCode != 200 {
		return "Couldn't reach API: " + resp.Status
	}

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

func Run(bot *irc.Connection, p, cmd, channel string, word []string) {
	say := bot.Privmsg

	if word[0] == p+"cat" || word[0] == p+"catfact" {
		say(channel, fact())
	}

	return
}
