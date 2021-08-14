package snd

import (
	"git.randomchars.net/freenitori/log"
	"os"
	"os/exec"
	"strconv"
)

var (
	channels         = 1
	channelsString   string
	sampleRate       int
	sampleRateString string
	frameSize        int
	maxSize          int
	ffmpegPath       string
)

var Conf struct {
	Mono       bool   `json:"mono"`
	SampleRate int    `json:"sample_rate"`
	FrameSize  int    `json:"frame_size"`
	FFmpegPath string `json:"ffmpeg_path"`
}

func Setup() error {
	ffmpegFind()

	if !Conf.Mono {
		channels++
	}

	sampleRate = Conf.SampleRate
	frameSize = Conf.FrameSize
	maxSize = frameSize * 4

	channelsString = strconv.Itoa(channels)
	sampleRateString = strconv.Itoa(sampleRate)

	return nil
}

func ffmpegFind() bool {
	if Conf.FFmpegPath == "auto" {
		if path, err := exec.LookPath("ffmpeg"); err != nil {
			log.Warn("FFmpeg was not found on this system.")
			return false
		} else {
			ffmpegPath = path
			log.Infof("FFmpeg discovered at %s.", ffmpegPath)
		}
	} else {
		if _, err := os.Stat(Conf.FFmpegPath); err != nil {
			log.Warn("Configured ffmpeg path does not exist, falling back to automatic lookup.")
			Conf.FFmpegPath = "auto"
			return ffmpegFind()
		} else {
			ffmpegPath = Conf.FFmpegPath
		}
	}
	return true
}
