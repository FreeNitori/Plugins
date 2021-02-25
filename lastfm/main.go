package main

import (
	"encoding/json"
	"errors"
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/log"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/shkh/lastfm-go/lastfm"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

const configPath = "plugins/lastfm.json"
const prefix = "LastFM: "

var err error
var conf struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// LastFM points to an instance of LastFM API client.
var LastFM *lastfm.Api

// Setup returns route after setting up LastFM session.
//goland:noinspection GoUnusedExportedFunction
func Setup() interface{} {
	if _, err := os.Stat(configPath); err != nil {
		log.Infof("%sNo config file found, generating default.", prefix)
		conf = struct {
			APIKey    string `json:"api_key"`
			APISecret string `json:"api_secret"`
		}(struct {
			APIKey    string
			APISecret string
		}{APIKey: "KEY_HERE", APISecret: "SECRET_HERE"})
		def, err := json.Marshal(conf)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(configPath, def, 0600)
		if err != nil {
			return err
		}
		log.Warnf("%sDefault config file generated, edit before next startup to enable LastFM.", prefix)
		return nil
	}
	log.Infof("%sLoading config file.", prefix)
	confData, err := ioutil.ReadFile("plugins/lastfm.json")
	if err != nil {
		return err
	}
	log.Infof("%sParsing config file.", prefix)
	err = json.Unmarshal(confData, &conf)
	if err != nil {
		return err
	}
	if conf.APIKey == "KEY_HERE" || conf.APISecret == "SECRET_HERE" {
		return errors.New("default configuration file was not edited")
	}
	LastFM = lastfm.New(conf.APIKey, conf.APISecret)
	return &multiplexer.Route{
		Pattern:       "fm",
		AliasPatterns: []string{"lastfm"},
		Description:   "Query last song scrobbled to lastfm.",
		Category:      multiplexer.MediaCategory,
		Handler:       fm,
	}
}

func fm(context *multiplexer.Context) {
	var username string
	switch len(context.Fields) {
	case 1:
	case 2:
		if context.Fields[1] == "unset" {
			err = resetLastfm(context.User, context.Guild)
			if !context.HandleError(err) {
				return
			}
			context.SendMessage("Successfully reset lastfm username.")
			return
		}
		username = context.Fields[1]
	case 3:
		switch context.Fields[1] {
		case "set":
			if b, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, context.Fields[2]); !b || len(context.Fields[2]) < 2 || len(context.Fields[2]) > 15 {
				context.SendMessage(multiplexer.InvalidArgument)
				return
			}
			err = setLastfm(context.User, context.Guild, context.Fields[2])
			if !context.HandleError(err) {
				return
			}
			context.SendMessage("Successfully set lastfm username to `" + context.Fields[2] + "`.")
			return
		case "unset":
			err = resetLastfm(context.User, context.Guild)
			if !context.HandleError(err) {
				return
			}
			context.SendMessage("Successfully reset lastfm username.")
			return
		default:
			context.SendMessage(multiplexer.InvalidArgument)
			return
		}
	}
	if username == "" {
		username, err = getLastfm(context.User, context.Guild)
	}
	if !context.HandleError(err) {
		return
	}
	p := lastfm.P{"user": username, "limit": 1, "extended": 0}
	result, err := LastFM.User.GetRecentTracks(p)
	if err != nil {
		context.SendMessage("Please set your lastfm username `" + context.Prefix() + "fm set <username>`.")
		return
	}
	if len(result.Tracks) < 1 {
		context.SendMessage("This username doesn't exist or does not have any scrobbles.")
		return
	}
	embed := embedutil.New(result.Tracks[0].Name, result.Tracks[0].Artist.Name+" | "+result.Tracks[0].Album.Name)
	embed.SetAuthor(context.User.Username, context.User.AvatarURL("128"))
	embed.SetFooter(fmt.Sprintf("%s has %s scrobbles in total.", result.User, strconv.Itoa(result.Total)))
	embed.Color = context.Session.State.UserColor(context.User.ID, context.Channel.ID)
	embed.URL = result.Tracks[0].Url
	if len(result.Tracks[0].Images) == 4 {
		embed.SetThumbnail(result.Tracks[0].Images[3].Url)
	}
	context.SendEmbed("", embed)
}
