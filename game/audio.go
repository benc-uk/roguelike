package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"path"
	"roguelike/core"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type soundEffects struct {
	audioContext *audio.Context
	sounds       map[string][]*audio.Player
	enabled      bool
}

func newSoundEffects(sfxPath string, enabled bool) *soundEffects {
	if !enabled {
		return &soundEffects{
			enabled:      false,
			sounds:       nil,
			audioContext: nil,
		}
	}

	es := &soundEffects{
		audioContext: audio.NewContext(44100),
		sounds:       make(map[string][]*audio.Player),
		enabled:      true,
	}

	es.loadSound("hurt", path.Join(sfxPath, "sfx"), 3)
	es.loadSound("walk", path.Join(sfxPath, "sfx"), 3)
	es.loadSound("pickup", path.Join(sfxPath, "sfx"), 1)
	return es
}

func (s *soundEffects) play(category string) {
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

func (s *soundEffects) loadSound(category string, soundPath string, count int) {
	s.sounds[category] = make([]*audio.Player, 0)

	for i := 0; i < count; i++ {
		d, err := core.ReadFile(path.Join(soundPath, fmt.Sprintf("%s_%d.wav", category, i)))
		if err != nil {
			log.Fatal(err)
		}

		player := s.audioContext.NewPlayerFromBytes(d)
		s.sounds[category] = append(s.sounds[category], player)
	}
}
