package Car

import (
	"strconv"
	"time"

	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

type Buzzer struct {
	driver *gpio.BuzzerDriver
}

func (b Buzzer) On() error {
	return b.driver.On()
}

func (b Buzzer) Off() error {
	return b.driver.Off()
}

func (b Buzzer) Toggle() error {
	return b.driver.Toggle()
}

func (b Buzzer) SetBPM(bpm float64) {
	b.driver.BPM = bpm
}

func (b Buzzer) PlayTone(note Note) error {
	if note.pitch <= 1 {
		time.Sleep(time.Duration(60 / b.driver.BPM * note.duration))
		return nil
	} else {
		return b.driver.Tone(note.pitch, note.duration)
	}
}

func (b Buzzer) PlaySong(song Song, finish chan struct{}, done chan error) error {
	for _, note := range song {
		select {
		case <-finish:
			{
				if done != nil {
					done <- nil
				}
				return nil
			}
		default:
			{
				err := b.PlayTone(note)
				if err != nil {
					if done != nil {
						done <- err
					}
					return err
				}
			}
		}
	}
	if done != nil {
		done <- nil
	}
	return nil
}

func CreateBuzzer() (Buzzer, error) {
	rpi := raspi.NewAdaptor()
	buzzer := gpio.NewBuzzerDriver(rpi, strconv.Itoa(17))
	err := buzzer.Start()
	if err != nil {
		return Buzzer{}, err
	}
	return Buzzer{buzzer}, nil
}
