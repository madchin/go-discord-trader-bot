package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func onReady(s *discordgo.Session, _ *discordgo.Ready) {
	log.Printf("websocket is ready for interaction \nbot user: %s", s.State.User.Username)
}

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	defer discord.Close()
	discord.AddHandler(onReady)
	if err := discord.Open(); err != nil {
		panic(err)
	}

}
