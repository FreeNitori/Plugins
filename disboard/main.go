package main

import (
	"encoding/json"
	"errors"
	"git.randomchars.net/freenitori/freenitori/v2/nitori/state"
	"git.randomchars.net/freenitori/log"
	"git.randomchars.net/freenitori/multiplexer"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const configPath = "plugins/disboard.json"
const prefix = "Disboard: "

var conf = struct {
	GuildIDs []int `json:"guild_ids"`
	UserID   int   `json:"user_id"`
}{}

var uid string
var gid map[string]bool

// Setup reads configuration file and sets up disboard parser and reminder.
//goland:noinspection GoUnusedExportedFunction
func Setup() interface{} {
	if _, err := os.Stat(configPath); err != nil {
		log.Infof("%sNo config file found, generating default.", prefix)
		conf = struct {
			GuildIDs []int `json:"guild_ids"`
			UserID   int   `json:"user_id"`
		}(struct {
			GuildIDs []int
			UserID   int
		}{GuildIDs: []int{}, UserID: 302050872383242240})
		def, err := json.Marshal(conf)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(configPath, def, 0600)
		if err != nil {
			return err
		}
		log.Warnf("%sDefault config file generated, edit before next startup to enable disboard parsing.", prefix)
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
	} else if conf.UserID == 0 {
		return errors.New("no bot user defined")
	}

	uid = strconv.Itoa(conf.UserID)
	gid = make(map[string]bool)
	for _, id := range conf.GuildIDs {
		gid[strconv.Itoa(id)] = true
	}

	state.Multiplexer.MessageCreate = append(state.Multiplexer.MessageCreate, disboardCreateHandler)
	return nil
}

func disboardCreateHandler(context *multiplexer.Context) {
	if context.Channel == nil || context.User == nil || context.Guild == nil || context.User.ID != uid ||
		len(context.Message.Embeds) != 1 || context.Message.Embeds[0] == nil || !gid[context.Guild.ID] {
		return
	}

	if strings.Contains(context.Message.Embeds[0].Description, "Please wait") {
		return
	}

	segments := strings.Split(context.Message.Embeds[0].Description, ",")
	bonker := context.GetMember(segments[0])
	if bonker == nil {
		return
	}

	if !context.HandleError(context.Session.ChannelTyping(context.Channel.ID)) {
		return
	}

	time.Sleep(1 * time.Second)
	context.SendMessage("<:FeelsKappa:779555772786802709>")
	go bonkTimer(context.Guild, bonker.User, context.SendMessage)
}

func bonkTimer(guild *discordgo.Guild, user *discordgo.User, messageSend func(message string) *discordgo.Message) {
	log.Infof("User %s bonked in jail %s.", user.ID, guild.ID)
	time.Sleep(120 * time.Minute)
	messageSend(user.Mention() + " BONK!")
}
