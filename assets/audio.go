package assets

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var (
	AudioContext *audio.Context
	sounds       map[string]*audio.Player
)

func LoadAudio() error {
	AudioContext = audio.NewContext(44100)
	sounds = make(map[string]*audio.Player)

	loadSound("hit", "assets/sounds/hit.wav")
	loadSound("shoot", "assets/sounds/shoot.wav")
	loadSound("death", "assets/sounds/death.wav")
	loadSound("restart", "assets/sounds/restart.wav")
	loadSound("jump", "assets/sounds/jump.wav")
	loadSound("music", "assets/sounds/music.wav")

	return nil
}

func loadSound(name, path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load sound %s: %v", path, err)
		return
	}
	s, err := wav.DecodeWithSampleRate(44100, bytes.NewReader(b))
	if err != nil {
		log.Printf("Failed to decode sound %s: %v", path, err)
		return
	}
	p, err := AudioContext.NewPlayer(s)
	if err != nil {
		log.Printf("Failed to create player %s: %v", path, err)
		return
	}
	sounds[name] = p
}

func PlaySound(name string) {
	if p, ok := sounds[name]; ok {
		p.Rewind()
		p.Play()
	}
}

func PlayMusic(name string) {
	if p, ok := sounds[name]; ok {
		if !p.IsPlaying() {
			p.Rewind()
			p.Play()
		}
	}
}
