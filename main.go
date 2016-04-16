package main

import (
	"encoding/json"
	"fmt"
	"github.com/thoj/go-ircevent"
	"log"
	"os"
)

type Bot struct {
	Nick                   string
	User                   string
	Server                 string
	Port                   int
	Channels               []string
	Debug                  bool
	VerboseCallbackHandler bool
	Admins                 map[string]string
	CommandPrefix          string
	irc                    *irc.Connection
}

func (bot *Bot) InitConnect() {
	log.Println("Connecting...")

	// connect
	bot.irc = irc.IRC(bot.Nick, bot.User)
	err := bot.irc.Connect(fmt.Sprintf("%s:%d", bot.Server, bot.Port))

	// log any errors
	if err != nil {
		log.Println("Error:", err)
	}

	// to help debugging, if set in config
	bot.irc.Debug = bot.Debug
	bot.irc.VerboseCallbackHandler = bot.VerboseCallbackHandler

	// join all channels in config
	for _, channel := range bot.Channels {
		log.Println(fmt.Sprintf("Joining %s", channel))
		bot.irc.Join(channel)
	}

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

	// open & readconfig
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)

	bot := Bot{}

	// decode the json
	err := decoder.Decode(&bot)
	if err != nil {
		log.Println("error: ", err)
	}

	bot.InitConnect()
}
