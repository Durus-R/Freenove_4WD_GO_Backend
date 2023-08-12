package Car

import (
	"errors"
	"log"
	"sync"
	"time"

	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

func RgbToColor(r int, g int, b int) uint32 {
	return uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}

type RGBStrip struct {
	ws2811 *ws281x.WS2811
	length int
	mutex  *sync.Mutex
	locked bool
}

func (r *RGBStrip) customColorWipe(color uint32, waitMs time.Duration) error {
	if r.locked {
		return errors.New("already locked")
	}
	r.locked = true
	r.mutex.Lock()
	for i := 0; i < r.length; i++ {

		r.ws2811.Leds(0)[i] = color
		err := r.ws2811.Render()
		if err != nil {
			r.mutex.Unlock()
			return err
		}
		time.Sleep(waitMs * time.Millisecond)
	}
	r.mutex.Unlock()
	return nil

}

func (r *RGBStrip) ColorWipe(color uint32) error {
	err := r.customColorWipe(color, 50)
	return err
}

func (r *RGBStrip) customTheaterChase(color uint32, waitMS time.Duration,
	iterations int, finish chan struct{}) error {
	if r.locked {
		return errors.New("already locked")
	}
	r.locked = true
	r.mutex.Lock()
	for i := 0; i < iterations; i++ {
		for j := 0; j < 3; j++ {
			select {
			case <-finish:
				{
					err := r.Black()
					r.locked = false
					r.mutex.Unlock()
					return err
				}
			default:
				for k := 0; k < r.length; k += 3 {
					r.ws2811.Leds(0)[k+j] = color
				}
				err := r.ws2811.Render()
				if err != nil {
					r.locked = false
					r.mutex.Unlock()
					return err
				}
				time.Sleep(waitMS * time.Millisecond)
				for k := 0; k < r.length; k += 3 {
					r.ws2811.Leds(0)[k+j] = 0
				}
			}

		}
	}
	r.locked = false
	r.mutex.Unlock()
	return nil
}

func (r *RGBStrip) TheaterChase(color uint32, finish chan struct{}) error {
	return r.customTheaterChase(color, 50, 10, finish)
}

func ColorWheel(pos uint32) uint32 {
	var r, g, b uint32
	if pos < 0 || pos > 255 {
		r, g, b = 0, 0, 0
	} else if pos < 85 {
		r = pos * 3
		g = 255 - pos*3
		b = 0
	} else if pos < 170 {
		pos -= 85
		r = 255 - pos*3
		g = 0
		b = pos * 3
	} else {
		pos -= 170
		r = 0
		g = pos * 3
		b = 255 - pos*3
	}
	return RgbToColor(int(r), int(g), int(b))
}

func (r *RGBStrip) customRainbow(waitMs time.Duration, iterations int,
	finish chan struct{}) error {
	if r.locked {
		return errors.New("already locked")
	}
	r.locked = true
	r.mutex.Lock()
	for i := 0; i < 256*iterations; i++ {
		select {
		case <-finish:
			{
				err := r.Black()
				r.locked = false
				r.mutex.Unlock()
				return err
			}
		default:
			{
				for j := 0; j < r.length; j++ {
					r.ws2811.Leds(0)[i] = ColorWheel(uint32(i+j) & 255)
				}
				err := r.ws2811.Render()
				if err != nil {
					r.locked = false
					r.mutex.Unlock()
					return err
				}
				time.Sleep(waitMs * time.Millisecond)
			}
		}

	}
	r.locked = false
	r.mutex.Unlock()
	return nil
}

func (r *RGBStrip) Rainbow(finish chan struct{}) error {
	return r.customRainbow(20, 1, finish)
}

func (r *RGBStrip) customRainbowCycle(waitMs time.Duration, iterations int,
	finish chan struct{}) error {
	if r.locked {
		return errors.New("already locked")
	}
	r.locked = true
	r.mutex.Lock()
	for i := 0; i < 256*iterations; i++ {
		select {
		case <-finish:
			{
				err := r.Black()
				r.locked = false
				r.mutex.Unlock()
				return err
			}
		default:
			{
				for j := 0; j < r.length; j++ {
					r.ws2811.Leds(0)[i] = ColorWheel(uint32(i*256/r.length + j&255))
				}
				err := r.ws2811.Render()
				if err != nil {
					r.locked = false
					r.mutex.Unlock()
					return err
				}
				time.Sleep(waitMs * time.Millisecond)
			}
		}
	}
	r.locked = false
	r.mutex.Unlock()
	return nil
}

func (r *RGBStrip) RainbowCycle(finish chan struct{}) error {
	return r.customRainbowCycle(20, 5, finish)
}

func (r *RGBStrip) customTheaterChaseRainbow(waitMs time.Duration,
	finish chan struct{}) error {
	if r.locked {
		return errors.New("already locked")
	}
	r.locked = true
	r.mutex.Lock()
	for i := 0; i < 256; i++ {
		for j := 0; j < 3; j++ {
			select {
			case <-finish:
				{
					err := r.Black()
					r.locked = false
					r.mutex.Unlock()
					return err
				}
			default:
				{
					for k := 0; k < r.length; k += 3 {
						r.ws2811.Leds(0)[k+j] = ColorWheel(uint32(k+j)) % 255
					}
					err := r.ws2811.Render()
					if err != nil {
						r.locked = false
						r.mutex.Unlock()
						return err
					}
					time.Sleep(waitMs * time.Millisecond)
					for k := 0; k < r.length; k += 3 {
						r.ws2811.Leds(0)[k+j] = 0
					}
				}
			}
		}
	}
	r.locked = false
	r.mutex.Unlock()
	return nil

}

func (r *RGBStrip) TheaterChaseRainbow(finish chan struct{}) error {
	return r.customTheaterChaseRainbow(50, finish)
}

func (r *RGBStrip) Black() error {
	return r.ColorWipe(RgbToColor(0, 0, 0))
}

func (r *RGBStrip) ApplyColors(c [8]uint32) error {
	if r.locked {
		return errors.New("already locked")
	}
	r.locked = true
	r.mutex.Lock()
	for i := range c {
		r.ws2811.Leds(0)[i] = c[i]
	}
	err := r.ws2811.Render()
	r.locked = false
	r.mutex.Unlock()
	return err
}

func (r *RGBStrip) IsLocked() bool {
	return r.locked
}

func CreateRGBStrip() RGBStrip {
	opt := ws281x.DefaultOptions
	opt.Channels[0].Brightness = 255
	opt.Channels[0].LedCount = 8
	opt.Channels[0].GpioPin = 18
	opt.Frequency = 800000

	ws2811, err := ws281x.MakeWS2811(&opt)
	if err != nil {
		log.Fatal("Error creating LED Strip: ", err)
	}
	err = ws2811.Init()
	if err != nil {
		log.Fatal("Error creating LED Strip: ", err)
	}
	defer ws2811.Fini()

	return RGBStrip{
		ws2811: ws2811,
		length: 8,
		mutex:  &sync.Mutex{},
	}

}
