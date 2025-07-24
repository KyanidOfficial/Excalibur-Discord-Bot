package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env only when not running in CI/CD (like GitHub Actions)
	if os.Getenv("CI") != "true" {
		_ = godotenv.Load()
	}

	// Get the token from environment variable (e.g., GitHub Secret)
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	if BOT_TOKEN == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	// Create Discord session
	sess, err := discordgo.New("Bot " + BOT_TOKEN)
	if err != nil {
		log.Fatal("Failed to create Discord session:", err)
	}

	// Register event handlers
	sess.AddHandler(onGuildCreate)
	sess.AddHandler(LeaveEveryServer)

	// Set intents
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open connection
	err = sess.Open()
	if err != nil {
		log.Fatal("Failed to open Discord session:", err)
	}
	defer sess.Close()

	// Set bot status
	_ = sess.UpdateStreamingStatus(0, "Excalibur / Blood Group", "https://www.twitch.tv/404")

	fmt.Println("Bot is online!")
	fmt.Println("[/] TOKEN: [REDACTED]")
	fmt.Printf("[/] LINK: https://discord.com/api/oauth2/authorize?client_id=%s&permissions=8&scope=bot\n", sess.State.User.ID)

	// Keep bot running
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}
