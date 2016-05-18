package common

import (
	"time"

	"github.com/mferrera/go-ircevent"
)

type Bot struct {
	Nick                   string            `json:"nick"`
	User                   string            `json:"user"`
	NickPW                 string            `json:"nickPw"`
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
	Irc                    *irc.Connection
}

type EmitMsg struct {
	P          string
	Cmd        string
	Word       []string
	Channel    string
	Channels   *[]string
	Nick       string
	Admin      bool
	Irc        *irc.Connection
	Say        func(string, string)
	IgnoreList map[string]string
}
