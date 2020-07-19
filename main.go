// Runs a Discord bot that listens for a couple simple commands: "ping" and "gif".
// Requires the following env variables:
// - "DISCORD_BOT_TOKEN": Token generate for the bot.
// - "ENVIRONMENT": Environment variable. Currently it just accepts "production" - any other value is assumed to be development.
// - "GIPHY_API_KEY": API key for the Giphy API.
package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joomcode/errorx"
	"go.uber.org/zap"
)

func initLogger(environment string) *zap.Logger {
	var logger *zap.Logger
	if IsProduction(environment) {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	zap.ReplaceGlobals(logger)
	return logger
}

func main() {
	environment := GetEnvironment()
	discordToken := GetDiscordBotToken()
	giphyAPIKey := GetGiphyAPIKey()

	logger := initLogger(environment)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		zap.S().Panic(errorx.Decorate(err, "Error while creating Discord session"))
	}

	giphy := NewGiphyHandler(giphyAPIKey)

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself.
		if m.Author.ID == s.State.User.ID {
			return
		}

		message := strings.TrimSpace(m.Content)
		zap.S().Debugf("Got message: %s", message)

		tokens := strings.Split(message, " ")
		if len(tokens) == 0 {
			zap.S().Debug("Found no tokens in message")
			return
		}

		command := tokens[0]
		commandParams := strings.Join(tokens[1:], " ")
		switch command {
		case "ping":
			s.ChannelMessageSend(m.ChannelID, "Pong!")

		case "gif":
			gif, err := giphy.SearchFirst(commandParams)
			if err != nil {
				zap.S().Error(errorx.Decorate(err, "Failed to fetch gif from Giphy"))
				return
			}

			s.ChannelMessageSend(m.ChannelID, gif.Url)
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		zap.S().Panic(errorx.Decorate(err, "Error while opening connection"))
	}

	// Wait here until CTRL-C or other term signal is received.
	zap.S().Info("Bot is now running. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close the Discord session before exiting.
	dg.Close()
	// Flush the logs before exiting.
	logger.Sync()
}
