package main

import (
	"git.randomchars.net/freenitori/freenitori/v2/nitori/database"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger/v3"
)

// getLastfm gets a user's lastfm username.
func getLastfm(user *discordgo.User, guild *discordgo.Guild) (string, error) {
	result, err := database.Database.HGet("lastfm."+guild.ID, user.ID)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return "", nil
		}
		return "", err
	}
	return result, err
}

// setLastfm sets a user's lastfm username.
func setLastfm(user *discordgo.User, guild *discordgo.Guild, username string) error {
	return database.Database.HSet("lastfm."+guild.ID, user.ID, username)
}

// resetLastfm resets a user's lastfm username.
func resetLastfm(user *discordgo.User, guild *discordgo.Guild) error {
	err := database.Database.HDel("lastfm."+guild.ID, []string{user.ID})
	if err == badger.ErrKeyNotFound {
		return nil
	}
	return err
}
