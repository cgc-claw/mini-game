package assets

import (
	"encoding/binary"
	"os"
)

func CreateSounds() error {
	os.MkdirAll("assets/sounds", 0755)

	if err := createHitSound(); err != nil {
		return err
	}
	if err := createShootSound(); err != nil {
		return err
	}
	if err := createDeathSound(); err != nil {
		return err
	}
	if err := createJumpSound(); err != nil {
		return err
	}
	if err := createRestartSound(); err != nil {
		return err
	}
	if err := createBackgroundMusic(); err != nil {
		return err
	}

	return nil
}

func createHitSound() error {
	return createTone("assets/sounds/hit.wav", 800, 100)
}

func createShootSound() error {
	return createTone("assets/sounds/shoot.wav", 1200, 80)
}

func createDeathSound() error {
	return createTone("assets/sounds/death.wav", 400, 300)
}

func createJumpSound() error {
	return createTone("assets/sounds/jump.wav", 1500, 100)
}

func createRestartSound() error {
	return createTone("assets/sounds/restart.wav", 1000, 150)
}

func createBackgroundMusic() error {
	return createTone("assets/sounds/music.wav", 600, 2000)
}

func createTone(filename string, frequency, duration int) error {
	sampleRate := 44100
	totalSamples := (sampleRate * duration) / 1000

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write WAV header
	writeWAVHeader(f, uint32(totalSamples))

	// Generate tone samples
	for i := 0; i < totalSamples; i++ {
		sample := int16(32767 * (i % (sampleRate / frequency)) / (sampleRate / frequency))
		if i%2 == 0 {
			sample = -sample
		}
		binary.Write(f, binary.LittleEndian, sample)
	}

	return nil
}

func writeWAVHeader(f *os.File, numSamples uint32) {
	sampleRate := uint32(44100)
	bytesPerSample := uint32(2)
	numChannels := uint16(1)
	byteRate := sampleRate * uint32(numChannels) * bytesPerSample
	blockAlign := uint16(numChannels * uint16(bytesPerSample))
	bitsPerSample := uint16(16)

	dataSize := numSamples * bytesPerSample * uint32(numChannels)

	// RIFF header
	f.Write([]byte("RIFF"))
	binary.Write(f, binary.LittleEndian, uint32(36+dataSize))
	f.Write([]byte("WAVE"))

	// fmt sub-chunk
	f.Write([]byte("fmt "))
	binary.Write(f, binary.LittleEndian, uint32(16))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, numChannels)
	binary.Write(f, binary.LittleEndian, sampleRate)
	binary.Write(f, binary.LittleEndian, byteRate)
	binary.Write(f, binary.LittleEndian, blockAlign)
	binary.Write(f, binary.LittleEndian, bitsPerSample)

	// data sub-chunk
	f.Write([]byte("data"))
	binary.Write(f, binary.LittleEndian, dataSize)
}
