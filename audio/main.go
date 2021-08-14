package main

import (
	"encoding/json"
	"git.randomchars.net/freenitori/log"
	"git.randomchars.net/freenitori/plugins/audio/snd"
	"io/ioutil"
	"os"
)

const configPath = "plugins/audio.json"
const prefix = "Audio: "

// Setup sets up audio-related stuff directly.
//goland:noinspection GoUnusedExportedFunction
func Setup() interface{} {
	if _, err := os.Stat(configPath); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		log.Infof("%sNo config file found, generating default.", prefix)
		snd.Conf = struct {
			Mono       bool   `json:"mono"`
			SampleRate int    `json:"sample_rate"`
			FrameSize  int    `json:"frame_size"`
			FFmpegPath string `json:"ffmpeg_path"`
		}(struct {
			Mono       bool
			SampleRate int
			FrameSize  int
			FFmpegPath string
		}{Mono: false, SampleRate: 48000, FrameSize: 960, FFmpegPath: "auto"})

		var def []byte
		if def, err = json.Marshal(snd.Conf); err != nil {
			return err
		}

		if err = ioutil.WriteFile(configPath, def, 0600); err != nil {
			return err
		}
	}

	log.Infof("%sLoading configuration.", prefix)
	if confData, err := ioutil.ReadFile(configPath); err != nil {
		return err
	} else {
		if err = json.Unmarshal(confData, &snd.Conf); err != nil {
			return err
		}
	}

	if err := snd.Setup(); err != nil {
		return err
	}

	return nil
}
