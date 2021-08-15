package main

import (
	"fmt"
	"git.randomchars.net/freenitori/freenitori/v2/nitori"
	"git.randomchars.net/freenitori/multiplexer"
	"git.randomchars.net/freenitori/plugins/audio/snd"
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

var s *snd.AudioSession

func play(context *multiplexer.Context) {
	path := "/usr/home/rand/bin/snd/Thousand Leaves/2013.05.25 [TIBA013] BLACK CLOAK SKELETON [Reitaisai 10]/08. Faithful Secret (Instrumental).flac"

	if s != nil {
		return
	}

	s = &snd.AudioSession{TextChannel: context.Channel}
	defer func() { s = nil }()

	if c, err := context.Session.ChannelVoiceJoin("713624993979695125", "762652462422687744", false, true); !context.HandleError(err) {
		return
	} else {
		defer func() {
			if !context.HandleError(c.Disconnect()) {
				return
			}
		}()
		s.VoiceConnection = c
	}

	context.SendMessage(fmt.Sprintf("Successfully connected to channel %s.", s.ChannelID))

	if !context.HandleError(s.Play(path)) {
		return
	}
}

func abort(context *multiplexer.Context) {
	if s == nil {
		return
	}

	if s.Abort() {
		context.SendMessage("Playback aborted.")
	} else {
		context.SendMessage("Uhhuh. Abort received with no abort channel.")
		context.SendMessage("Do you have a resource leak?")
		context.SendMessage("Dazed and confused, but trying to continue.")
	}
}
