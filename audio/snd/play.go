package snd

import (
	"bufio"
	"encoding/binary"
	"errors"
	"git.randomchars.net/freenitori/log"
	"github.com/bwmarrin/discordgo"
	"io"
	"layeh.com/gopus"
	"os/exec"
)

var ErrNoFFmpeg = errors.New("no ffmpeg found on this system")
var ErrInvalidArgument = errors.New("invalid argument")

// Out sends received PCM data to a discordgo.VoiceConnection.
func Out(connection *discordgo.VoiceConnection, pcm <-chan []int16) error {
	if pcm == nil {
		return ErrInvalidArgument
	}

	var encoder *gopus.Encoder
	if e, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio); err != nil {
		return err
	} else {
		encoder = e
	}

	for {
		p, ok := <-pcm
		if !ok {
			log.Warn("Audio out channel close.")
			return ErrInvalidArgument
		}

		o, err := encoder.Encode(p, frameSize, maxSize)
		if err != nil {
			return err
		}

		if connection.Ready == false || connection.OpusSend == nil {
			log.Debug("Audio out connection no longer ready.")
			return nil
		}

		connection.OpusSend <- o
	}
}

// PlayOnConnection plays audio on a discordgo.VoiceConnection.
func PlayOnConnection(connection *discordgo.VoiceConnection, path string, abort <-chan bool) error {
	if ffmpegPath == "" {
		if !ffmpegFind() {
			return ErrNoFFmpeg
		}
	}

	var stdout io.ReadCloser
	var running = true

	ffmpeg := exec.Command(ffmpegPath,
		"-i", path, "-f", "s16le", ""+"-ar", sampleRateString, "-ac", channelsString, "pipe:1")

	if out, err := ffmpeg.StdoutPipe(); err != nil {
		return err
	} else {
		stdout = out
	}

	defer func() {
		_ = stdout.Close()
	}()

	ffmpegOutBuf := bufio.NewReaderSize(stdout, 16384)

	if err := ffmpeg.Start(); err != nil {
		return err
	} else {
		go func() {
			err = ffmpeg.Wait()
			running = false
			if err != nil {
				log.Warnf("FFmpeg exited abnormally, %s", err)
			} else {
				log.Debugf("FFmpeg PID %v exited normally.", ffmpeg.Process.Pid)
			}
		}()
	}

	defer func() {
		if running {
			if err := ffmpeg.Process.Kill(); err != nil {
				log.Warnf("Error killing ffmpeg, %s", err)
			}
		}

		if err := connection.Speaking(false); err != nil {
			log.Warnf("Error sending speaking stop notification on channel %s, %s",
				connection.ChannelID, err)
		}
	}()

	go func() {
		if <-abort {
			if running {
				if err := ffmpeg.Process.Kill(); err != nil {
					log.Warnf("Error killing ffmpeg, %s", err)
				}
			}
		}
	}()

	if err := connection.Speaking(true); err != nil {
		return err
	}

	out := make(chan []int16, 2)
	defer close(out)

	outErr := make(chan error)
	go func() {
		outErr <- Out(connection, out)
	}()

	for {
		audioBuf := make([]int16, frameSize*channels)
		if err := binary.Read(ffmpegOutBuf, binary.LittleEndian, &audioBuf); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Debugf("ffmpeg buffer EOF")
				return nil
			}
			return err
		} else {
			select {
			case out <- audioBuf:
			case err = <-outErr:
				return err
			}
		}
	}
}
