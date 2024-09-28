package audio

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"path"
	"roguelike/core"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"gopkg.in/yaml.v3"
)

type Effects struct {
	audioContext *audio.Context
	sounds       map[string][]*audio.Player
	enabled      bool
}

type soundsMetaFile struct {
	SourceDir string         `yaml:"sourceDir"`
	Effects   map[string]int `yaml:"effects"`
}

func NewEffects(metaFile string, enabled bool) (*Effects, error) {
	// This allows for creation, but disables and stubs it out, to disable sound entirely
	if !enabled {
		return &Effects{
			enabled:      false,
			sounds:       nil,
			audioContext: nil,
		}, nil
	}

	data, err := core.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}

	var soundsFile soundsMetaFile
	err = yaml.Unmarshal(data, &soundsFile)
	if err != nil {
		return nil, err
	}

	efx := &Effects{
		audioContext: audio.NewContext(44100),
		sounds:       make(map[string][]*audio.Player),
		enabled:      true,
	}

	metaDir := path.Dir(metaFile)

	totalSounds := 0
	for category, count := range soundsFile.Effects {
		totalSounds += count

		err := efx.loadWavSet(path.Join(metaDir, soundsFile.SourceDir), category, count)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("Loaded %d sound effects", totalSounds)

	return efx, nil
}

func (s *Effects) Play(category string) {
	if !s.enabled {
		return
	}

	_, ok := s.sounds[category]
	if !ok {
		return
	}

	playerSlice := s.sounds[category]
	if len(playerSlice) == 0 {
		return
	}

	// Play a random sound in the category
	i := rand.IntN(len(playerSlice))
	_ = playerSlice[i].Rewind()
	playerSlice[i].Play()
}

func (s *Effects) loadWavSet(wavDir string, category string, count int) error {
	s.sounds[category] = make([]*audio.Player, 0)

	for i := 0; i < count; i++ {
		wavFile := path.Join(wavDir, fmt.Sprintf("%s_%d.wav", category, i))
		fileRawBytes, err := core.ReadFile(wavFile)
		if err != nil {
			return err
		}

		stream, err := wav.DecodeWithSampleRate(44100, bytes.NewReader(fileRawBytes))
		if err != nil {
			return err
		}

		wavBytes, err := io.ReadAll(stream)
		if err != nil {
			return err
		}

		player := s.audioContext.NewPlayerFromBytes(wavBytes)
		s.sounds[category] = append(s.sounds[category], player)
	}

	return nil
}
