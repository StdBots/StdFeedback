package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration.
type Config struct {
	BotToken        string
	OwnerID         int64
	AuthUsers       []int64
	LogChannel      int64
	DBURL           string
	DBName          string
	MongoURI        string
	MongoDBName     string
	StartText       string
	DonateLink      string
	BroadcastAsCopy bool
}

// Load reads configuration from environment variables.
func Load() *Config {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Panic("BOT_TOKEN is required")
	}

	ownerIDStr := os.Getenv("OWNER_ID")
	if ownerIDStr == "" {
		log.Panic("OWNER_ID is required")
	}
	ownerID, err := strconv.ParseInt(ownerIDStr, 10, 64)
	if err != nil {
		log.Panicf("Invalid OWNER_ID: %v", err)
	}

	var authUsers []int64
	authUsersStr := os.Getenv("AUTH_USERS")
	if authUsersStr != "" {
		for _, u := range strings.Fields(authUsersStr) {
			id, err := strconv.ParseInt(u, 10, 64)
			if err == nil {
				authUsers = append(authUsers, id)
			}
		}
	}

	logChannelStr := os.Getenv("LOG_CHANNEL")
	var logChannel int64
	if logChannelStr != "" {
		logChannel, _ = strconv.ParseInt(logChannelStr, 10, 64)
	}

	dbURL := os.Getenv("MONGO_URI")
	if dbURL == "" {
		dbURL = os.Getenv("DB_URL")
	}
	if dbURL == "" {
		log.Panic("MONGO_URI or DB_URL is required")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "stdfeedback"
	}

	startText := os.Getenv("START_TEXT")
	donateLink := os.Getenv("DONATE_LINK")

	broadcastAsCopy := true
	if strings.ToLower(os.Getenv("BROADCAST_AS_COPY")) == "false" {
		broadcastAsCopy = false
	}

	return &Config{
		BotToken:        botToken,
		OwnerID:         ownerID,
		AuthUsers:       authUsers,
		LogChannel:      logChannel,
		DBURL:           dbURL,
		DBName:          dbName,
		MongoURI:        dbURL,
		MongoDBName:     dbName,
		StartText:       startText,
		DonateLink:      donateLink,
		BroadcastAsCopy: broadcastAsCopy,
	}
}

// MustLoad is an alias to Load.
func MustLoad() *Config {
	return Load()
}
