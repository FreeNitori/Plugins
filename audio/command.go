package main

import (
	"fmt"
	"git.randomchars.net/freenitori/freenitori/v2/nitori"
	"git.randomchars.net/freenitori/multiplexer"
	"git.randomchars.net/freenitori/plugins/audio/snd"
	"github.com/bwmarrin/discordgo"
)

func init() {
	nitori.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "play",
		AliasPatterns: []string{"p"},
		Description:   "Play audio in a voice channel.",
		Category:      multiplexer.AudioCategory,
		Handler:       play,
	})

	nitori.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "abort",
		AliasPatterns: []string{"s", "stop"},
		Description:   "Abort audio playback.",
		Category:      multiplexer.AudioCategory,
		Handler:       abort,
	})
}

var a chan bool

func play(context *multiplexer.Context) {
	path := "/usr/home/rand/bin/snd/Thousand Leaves/2013.05.25 [TIBA013] BLACK CLOAK SKELETON [Reitaisai 10]/08. Faithful Secret (Instrumental).flac"

	if a != nil {
		return
	}

	var connection *discordgo.VoiceConnection
	if c, err := context.Session.ChannelVoiceJoin("713624993979695125", "762652462422687744", false, true); !context.HandleError(err) {
		return
	} else {
		defer func() {
			if !context.HandleError(c.Disconnect()) {
				return
			}
		}()
		connection = c
	}

	context.SendMessage(fmt.Sprintf("Successfully connected to channel %s.", connection.ChannelID))

	a = make(chan bool)
	defer func() { a = nil }()
	if !context.HandleError(snd.PlayOnConnection(connection, path, a)) {
		return
	}
}

func abort(context *multiplexer.Context) {
	if a == nil {
		return
	}

	a <- true
	context.SendMessage("Playback aborted.")
}
