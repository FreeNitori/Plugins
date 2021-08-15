package snd

import (
	"bufio"
	"encoding/binary"
	"errors"
	"git.randomchars.net/freenitori/log"
	"io"
	"os/exec"
)

var ErrAbortSet = errors.New("abort is still set for a new playback")

// Play plays an audio file.
func (s *AudioSession) Play(path string) error {
	s.Lock()

	if s.abort != nil {
		return ErrAbortSet
	} else {
		s.abort = make(chan bool)
	}

	defer func() {
		s.abort = nil
		s.Unlock()
	}()

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

	ffmpegOutBuf := bufio.NewReaderSize(stdout, 16384)

	if err := ffmpeg.Start(); err != nil {
		return err
	} else {
		go func() {
			err = ffmpeg.Wait()
			running = false
			if err != nil {
				if err.Error() == "signal: killed" {
					return
				}
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

		if err := s.Speaking(false); err != nil {
			log.Warnf("Error sending speaking stop notification on channel %s, %s",
				s.ChannelID, err)
		}
	}()

	go func() {
		if <-s.abort {
			if running {
				if err := ffmpeg.Process.Kill(); err != nil {
					log.Warnf("Error killing ffmpeg, %s", err)
				}
			}
		}
	}()

	if err := s.Speaking(true); err != nil {
		return err
	}

	out := make(chan []int16, 2)
	defer close(out)

	outErr := make(chan error)
	go func() {
		outErr <- s.Out(out)
	}()

	for running {
		audioBuf := make([]int16, frameSize*channels)
		if err := binary.Read(ffmpegOutBuf, binary.LittleEndian, &audioBuf); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Debug("ffmpeg buffer EOF")
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

	return nil
}

// Abort aborts audio playback.
func (s *AudioSession) Abort() bool {
	if s.abort == nil {
		return false
	}
	s.abort <- true
	return true
}
