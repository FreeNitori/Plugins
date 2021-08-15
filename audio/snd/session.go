package snd

import (
	"errors"
	"git.randomchars.net/freenitori/log"
	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
	"sync"
)

var ErrNoFFmpeg = errors.New("no ffmpeg found on this system")
var ErrInvalidArgument = errors.New("invalid argument")

var decoder = make(map[uint32]*gopus.Decoder)

type AudioSession struct {
	TextChannel *discordgo.Channel
	*discordgo.VoiceConnection
	abort chan bool
	sync.Mutex
}

// Out sends received PCM data.
func (s *AudioSession) Out(pcm <-chan []int16) error {
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
		if s.Ready == false || s.OpusSend == nil {
			log.Debug("Audio out connection no longer ready.")
			return nil
		}

		p, ok := <-pcm
		if !ok {
			log.Warn("Audio out channel close.")
			return ErrInvalidArgument
		}

		o, err := encoder.Encode(p, frameSize, maxSize)
		if err != nil {
			return err
		}

		s.OpusSend <- o
	}
}

// In receives opus and decodes it into PCM.
func (s *AudioSession) In(packet chan *discordgo.Packet) error {
	if packet == nil {
		return ErrInvalidArgument
	}

	for {
		if s.Ready == false || s.OpusRecv == nil {
			log.Debug("Audio in connection no longer ready.")
			return nil
		}

		var pkt *discordgo.Packet
		if p, ok := <-s.OpusRecv; !ok {
			log.Warn("Audio in channel close.")
			return ErrInvalidArgument
		} else {
			pkt = p
		}

		if _, ok := decoder[pkt.SSRC]; !ok {
			if d, err := gopus.NewDecoder(sampleRate, channels); err != nil {
				return err
			} else {
				decoder[pkt.SSRC] = d
			}
		}

		if pcm, err := decoder[pkt.SSRC].Decode(pkt.Opus, frameSize, false); err != nil {
			return err
		} else {
			pkt.PCM = pcm
		}

		packet <- pkt
	}
}
