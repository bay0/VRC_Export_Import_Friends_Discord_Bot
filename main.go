package main

import (
	"vrc_bot/bot"
	"vrc_bot/config"
	"vrc_bot/logging"
)

func init() {
	logging.Init()
	config.Init()
}

func main() {
	bot.Init()
}
