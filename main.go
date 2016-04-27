package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mferrera/go-ircevent"
)

type Bot struct {
	Nick                   string            `json:"nick"`
	User                   string            `json:"user"`
	Server                 string            `json:"server"`
	Port                   int               `json:"port"`
	Channels               []string          `json:"channels"`
	FloodProtect           bool              `json:"floodProtect"`
	FloodDelay             time.Duration     `json:"floodDelay"`
	Debug                  bool              `json:"debug"`
	VerboseCallbackHandler bool              `json:"verbose"`
	Admins                 map[string]string `json:"admins"`
	CommandPrefix          string            `json:"commandPrefix"`
	IgnoreList             map[string]string `json:"ignoreList"`
	DBAddress              string            `json:"dbAddress"`
	irc                    *irc.Connection
}

func (bot *Bot) connect() {
	log.Println("Connecting...")

	// connect
	bot.irc = irc.IRC(bot.Nick, bot.User)

	err := bot.irc.Connect(fmt.Sprintf("%s:%d", bot.Server, bot.Port))
	if err != nil {
		log.Panic("Error:", err)
	}

	// set values from conf.json to their irc.Connection equivalent
	bot.irc.FloodProtect = bot.FloodProtect
	bot.irc.FloodDelay = bot.FloodDelay
	bot.irc.Debug = bot.Debug
	bot.irc.VerboseCallbackHandler = bot.VerboseCallbackHandler
	bot.IgnoreList[bot.irc.GetNick()] = "Ignore messages from self."

	// join all channels in config
	bot.irc.AddCallback("001", func(e *irc.Event) {
		for _, channel := range bot.Channels {
			log.Println(fmt.Sprintf("Joining %s", channel))
			bot.irc.Join(channel)
		}
	})

	// log admins
	for admin, hostname := range bot.Admins {
		log.Println(fmt.Sprintf("Admin: %s!*@%s", admin, hostname))
	}

	// add listener callback
	listen(bot)

	// stay connected
	bot.irc.Loop()
}

func main() {
	log.Println("Starting up...")
	log.Println("Reading conf.json...")

	// open & read config
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal("Problem reading config file. " +
			"Make sure you renamed conf-example.json to conf.json " +
			"and properly edited it.")
	}
	decoder := json.NewDecoder(file)

	bot := Bot{}

	// decode the json
	err = decoder.Decode(&bot)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	bot.connect()
}
