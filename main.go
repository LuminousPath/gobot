package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/davecheney/profile"
	"github.com/mferrera/go-ircevent"
	"github.com/mferrera/gobot/common"
)

func connect(bot common.Bot) {
	log.Println("Connecting...")

	// connect
	bot.Irc = irc.IRC(bot.Nick, bot.User)

	err := bot.Irc.Connect(fmt.Sprintf("%s:%d", bot.Server, bot.Port))
	if err != nil {
		log.Panic("Error:", err)
	}

	// set values from conf.json to their Irc.Connection equivalent
	bot.Irc.FloodProtect = bot.FloodProtect
	bot.Irc.FloodDelay = bot.FloodDelay
	bot.Irc.Debug = bot.Debug
	bot.Irc.VerboseCallbackHandler = bot.VerboseCallbackHandler
	bot.IgnoreList[bot.Irc.GetNick()] = "Ignore messages from self."

	// join all channels in config
	bot.Irc.AddCallback("001", func(e *irc.Event) {
		for _, channel := range bot.Channels {
			log.Println(fmt.Sprintf("Joining %s", channel))
			bot.Irc.Join(channel)
		}
	})

	// log admins
	for admin, hostname := range bot.Admins {
		log.Println(fmt.Sprintf("Admin: %s!*@%s", admin, hostname))
	}

	if bot.NickPW != "" {
		bot.Irc.Privmsg("NickServ", "identify "+bot.NickPW)
	}

	// add listener callback
	listen(bot)

	// stay connected
	bot.Irc.Loop()
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	defer profile.Start(profile.MemProfile).Stop()
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

	bot := common.Bot{}

	// decode the json
	err = decoder.Decode(&bot)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	connect(bot)
}
