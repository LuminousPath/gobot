package catfact

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mferrera/gobot/common"
)

type CatFact struct {
	Fact    []string `json:"facts"`
	Success string   `json:"success"`
}

var (
	api  string = "http://catfacts-api.appspot.com/api/facts?number=1"
	cat  CatFact
	say  func(string, string)
	resp *http.Response
	err  error
)

func fact() string {
	resp, err = http.Get(api)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "Couldn't reach API: " + resp.Status
	}

	decoder := json.NewDecoder(resp.Body)

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

func Run(m common.EmitMsg) {
	if m.Word[0] == m.P+"cat" || m.Word[0] == m.P+"catfact" {
		m.Say(m.Channel, fact())
	}
}
