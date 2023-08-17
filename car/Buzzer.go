package car

import (
	"errors"
	"strconv"
	"time"

	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

type Buzzer struct {
	driver *gpio.BuzzerDriver
}

func (b *Buzzer) On() error {
	return b.driver.On()
}

func (b *Buzzer) Off() error {
	return b.driver.Off()
}

func (b *Buzzer) Toggle() error {
	return b.driver.Toggle()
}

func (b *Buzzer) SetBPM(bpm float64) error {
	if bpm < 1 {
		return errors.New("BPM is too low")
	}
	b.driver.BPM = bpm
	return nil
}

func (b *Buzzer) GetBPM() float64 {
	return b.driver.BPM
}

func (b *Buzzer) PlayTone(note Note) error {
	if note.Pitch <= 1 {
		time.Sleep(time.Duration(60 / b.driver.BPM * note.Duration))
		return nil
	} else {
		return b.driver.Tone(note.Pitch, note.Duration)
	}
}

func (b *Buzzer) PlaySong(song Song, finish chan struct{}) error {
	for _, note := range song {
		select {
		case <-finish:
			{
				return nil
			}
		default:
			{
				err := b.PlayTone(note)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func NewBuzzer() (*Buzzer, error) {
	rpi := raspi.NewAdaptor()
	buzzer := gpio.NewBuzzerDriver(rpi, strconv.Itoa(17))
	err := buzzer.Start()
	if err != nil {
		return &Buzzer{}, err
	}
	return &Buzzer{buzzer}, nil
}
