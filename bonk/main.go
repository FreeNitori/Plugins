package main

import (
	"encoding/json"
	"errors"
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"strconv"
)

const configPath = "plugins/bonk.json"
const prefix = "Bonk: "

var conf = struct {
	GuildIDs []int `json:"guild_ids"`
	UserIDs  []int `json:"user_ids"`
}{}

// Setup reads configuration file and sets up bonk tool.
//goland:noinspection GoUnusedExportedFunction
func Setup() interface{} {
	if _, err := os.Stat(configPath); err != nil {
		log.Infof("%sNo config file found, generating default.", prefix)
		conf = struct {
			GuildIDs []int `json:"guild_ids"`
			UserIDs  []int `json:"user_ids"`
		}(struct {
			GuildIDs []int
			UserIDs  []int
		}{GuildIDs: []int{}, UserIDs: []int{}})
		def, err := json.Marshal(conf)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(configPath, def, 0600)
		if err != nil {
			return err
		}
		log.Warnf("%sDefault config file generated, edit before next startup to enable bonking.", prefix)
		return nil
	}
	log.Infof("%sLoading config file.", prefix)
	confData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	log.Infof("%sParsing config file.", prefix)
	err = json.Unmarshal(confData, &conf)
	if err != nil {
		return err
	}
	if len(conf.GuildIDs) == 0 {
		return errors.New("no guilds are defined")
	} else if len(conf.UserIDs) == 0 {
		return errors.New("no users defined")
	}
	state.Multiplexer.MessageDelete = append(state.Multiplexer.MessageDelete, bonkDeleteHandler)
	state.Multiplexer.MessageUpdate = append(state.Multiplexer.MessageUpdate, bonkUpdateHandler)
	return nil
}

func bonkDeleteHandler(context *multiplexer.Context) {
	messageDelete, ok := context.Event.(*discordgo.MessageDelete)
	if !ok {
		return
	}
	if messageDelete.GuildID == "" {
		return
	}
	if messageDelete.BeforeDelete == nil {
		return
	}
	if !guildMatch(messageDelete.GuildID) || !userMatch(messageDelete.BeforeDelete.Author.ID) {
		return
	}
	var embed = embedutil.New("Message Delete", "")
	embed.Color = multiplexer.KappaColor
	embed.SetAuthor(messageDelete.BeforeDelete.Author.Username+"#"+messageDelete.BeforeDelete.Author.Discriminator, messageDelete.BeforeDelete.Author.AvatarURL("128"))
	embed.SetFooter(fmt.Sprintf("Channel: %s Message: %s Author: %s", messageDelete.ChannelID, messageDelete.BeforeDelete.ID, messageDelete.BeforeDelete.Author.ID))
	if messageDelete.BeforeDelete.Content != "" {
		embed.AddField("Content Pre", messageDelete.BeforeDelete.Content, false)
	}
	for _, attachment := range messageDelete.BeforeDelete.Attachments {
		embed.AddField("Attachment Pre", fmt.Sprintf("[%s](%s)", attachment.Filename, attachment.URL), false)
	}
	if messageDelete.BeforeDelete.MessageReference != nil {
		embed.AddField("References", fmt.Sprintf("[Message Link](https://discord.com/channels/%s/%s/%s)",
			messageDelete.BeforeDelete.MessageReference.GuildID,
			messageDelete.BeforeDelete.MessageReference.ChannelID,
			messageDelete.BeforeDelete.MessageReference.MessageID), false)
	}
	embed.AddField("Channel", fmt.Sprintf("<#%s>", messageDelete.ChannelID), false)
	context.SendEmbed("", embed)
	for _, e := range messageDelete.BeforeDelete.Embeds {
		context.SendEmbed("Embed included in previously deleted message.", embedutil.Embed{MessageEmbed: e})
	}
}

func bonkUpdateHandler(context *multiplexer.Context) {
	update, ok := context.Event.(*discordgo.MessageUpdate)
	if !ok {
		return
	}
	if update.GuildID == "" {
		return
	}
	if update.BeforeUpdate == nil {
		return
	}
	if update.Author == nil {
		return
	}
	if !guildMatch(update.GuildID) || !userMatch(update.Author.ID) {
		return
	}
	var embed = embedutil.New("Message Update",
		fmt.Sprintf("[Message Link](https://discord.com/channels/%s/%s/%s)",
			update.BeforeUpdate.GuildID,
			update.BeforeUpdate.ChannelID,
			update.BeforeUpdate.ID))
	embed.Color = multiplexer.KappaColor
	embed.SetAuthor(update.BeforeUpdate.Author.Username+"#"+update.BeforeUpdate.Author.Discriminator, update.BeforeUpdate.Author.AvatarURL("128"))
	embed.SetFooter(fmt.Sprintf("Channel: %s Message: %s Author: %s", update.ChannelID, update.BeforeUpdate.ID, update.BeforeUpdate.Author.ID))
	if update.BeforeUpdate.Content != "" {
		embed.AddField("Content Pre", update.BeforeUpdate.Content, false)
	}
	for _, attachment := range update.BeforeUpdate.Attachments {
		embed.AddField("Attachment Pre", fmt.Sprintf("[%s](%s)", attachment.Filename, attachment.URL), false)
	}
	if update.Message.Content != "" {
		embed.AddField("Content Post", update.Message.Content, false)
	}
	for _, attachment := range update.Message.Attachments {
		embed.AddField("Attachment Post", fmt.Sprintf("[%s](%s)", attachment.Filename, attachment.URL), false)
	}
	if update.BeforeUpdate.MessageReference != nil {
		embed.AddField("References", fmt.Sprintf("[Message Link](https://discord.com/channels/%s/%s/%s)",
			update.BeforeUpdate.MessageReference.GuildID,
			update.BeforeUpdate.MessageReference.ChannelID,
			update.BeforeUpdate.MessageReference.MessageID), false)
	}
	embed.AddField("Channel", fmt.Sprintf("<#%s>", update.ChannelID), false)
	context.SendEmbed("", embed)
	for _, e := range update.BeforeUpdate.Embeds {
		context.SendEmbed("Embed included in previously updated message.", embedutil.Embed{MessageEmbed: e})
	}
}

func guildMatch(tid string) bool {
	for _, id := range conf.GuildIDs {
		if strconv.Itoa(id) == tid {
			return true
		}
	}
	return false
}

func userMatch(tid string) bool {
	for _, id := range conf.UserIDs {
		if strconv.Itoa(id) == tid {
			return true
		}
	}
	return false
}
