package Car

import (
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
	"time"
)

type Ultrasonic struct {
	trigPin *gpio.DirectPinDriver
	echoPin *gpio.DirectPinDriver
}

func (u *Ultrasonic) Probe() (float64, error) {
	err := u.trigPin.DigitalWrite(byte(0))
	if err != nil {
		return 0, err
	}
	time.Sleep(2 * time.Microsecond)

	err = u.trigPin.DigitalWrite(byte(1))
	if err != nil {
		return 0, err
	}
	time.Sleep(10 * time.Microsecond)

	err = u.trigPin.DigitalWrite(byte(0))
	if err != nil {
		return 0, err
	}
	start := time.Now()
	end := time.Now()

	for {
		val, err := u.echoPin.DigitalRead()
		start = time.Now()

		if err != nil {
			return 0, err
		}

		if val == 0 {
			continue
		}

		break
	}

	for {
		val, err := u.echoPin.DigitalRead()
		end = time.Now()
		if err != nil {
			return 0, err
		}

		if val == 1 {
			continue
		}

		break
	}

	duration := end.Sub(start)
	distance := duration.Seconds() * 34300
	distance = distance / 2 //one way travel time
	return distance, nil
}

func NewUltrasonic() *Ultrasonic {
	adaptor := raspi.NewAdaptor()
	trigger := gpio.NewDirectPinDriver(adaptor, "27")
	echo := gpio.NewDirectPinDriver(adaptor, "22")
	return &Ultrasonic{
		trigPin: trigger,
		echoPin: echo,
	}
}
