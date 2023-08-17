package car

import (
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/raspi"
	"log"
	"math"
)

type MotorController struct {
	pca9685 *i2c.PCA9685Driver
}

type ChannelMap struct {
	positiveChannel int
	negativeChannel int
}

type Direction struct {
	LeftUp   int
	LeftLow  int
	RightUp  int
	RightLow int
}

func GetDirectionForward() Direction {
	return Direction{
		LeftUp:   2000,
		LeftLow:  2000,
		RightUp:  2000,
		RightLow: 2000,
	}
}

func GetDirectionBackward() Direction {
	return Direction{
		LeftUp:   -2000,
		LeftLow:  -2000,
		RightUp:  -2000,
		RightLow: -2000,
	}
}

func GetDirectionLeft() Direction {
	return Direction{
		LeftUp:   -500,
		LeftLow:  -500,
		RightUp:  2000,
		RightLow: 2000,
	}
}

func GetDirectionRight() Direction {
	return Direction{
		LeftUp:   2000,
		LeftLow:  2000,
		RightUp:  -500,
		RightLow: -500,
	}
}

func GetDirectionHalt() Direction {
	return Direction{
		LeftUp:   0,
		LeftLow:  0,
		RightUp:  0,
		RightLow: 0,
	}
}

func GetChannelMapLeftUp() ChannelMap {
	return ChannelMap{0, 1}
}

func GetChannelMapLeftLow() ChannelMap {
	return ChannelMap{3, 2}
}

func GetChannelMapRightUp() ChannelMap {
	return ChannelMap{6, 7}
}

func GetChannelMapRightLow() ChannelMap {
	return ChannelMap{4, 5}
}

func NewMotorController() (*MotorController, error) {
	adaptor := raspi.NewAdaptor()
	pca9685 := i2c.NewPCA9685Driver(adaptor, i2c.WithBus(1), i2c.WithAddress(0x40))

	err := pca9685.Start()
	if err != nil {
		return nil, err
	}

	err = pca9685.SetPWMFreq(50.0)
	if err != nil {
		log.Print("Could not set Frequency: ", err)
		return nil, err
	}

	return &MotorController{pca9685: pca9685}, nil
}

func (m *MotorController) SetAngle(channel int, angle uint16) {
	res := angle + 111
	switch channel {
	case 0:
		{
			err := m.pca9685.SetPWM(8, 0, 2500-res)
			if err != nil {
				log.Fatal("Error on setting Servo: ", err)
			}
		}
	case 1:
		{
			err := m.pca9685.SetPWM(9, 0, 500+res)
			if err != nil {
				log.Fatal("Error on setting Servo: ", err)
			}
		}

	}

}

func (m *MotorController) SetMotorDuty(c ChannelMap, duty int) {
	if duty == 0 {
		m.SetNullDuty(c)
		return
	} else if duty > 4095 {
		duty = 4095
	} else if duty < -4095 {
		duty = -4095
	}
	m.SetCorrectedDuty(c, duty)
}

func (m *MotorController) SetCorrectedDuty(c ChannelMap, duty int) {
	if duty < 0 {
		err := m.pca9685.SetPWM(c.negativeChannel, 0, 0)
		if err != nil {
			log.Print("Error in PWM: ", err)
		}
		err = m.pca9685.SetPWM(c.positiveChannel, 0, uint16(math.Abs(float64(duty))))
		if err != nil {
			log.Print("Error in PWM: ", err)
		}
	} else {
		err := m.pca9685.SetPWM(c.positiveChannel, 0, 0)
		if err != nil {
			log.Print("Error in PWM: ", err)
		}
		err = m.pca9685.SetPWM(c.negativeChannel, 0, uint16(math.Abs(float64(duty))))
		if err != nil {
			log.Print("Error in PWM: ", err)
		}
	}
}

func (m *MotorController) SetNullDuty(c ChannelMap) {
	err := m.pca9685.SetPWM(c.negativeChannel, 0, 4095)
	if err != nil {
		log.Print("Error in PWM: ", err)
	}
	err = m.pca9685.SetPWM(c.positiveChannel, 0, 4095)
	if err != nil {
		log.Print("Error in PWM: ", err)
	}

}

func (m *MotorController) SetDirection(d Direction) {
	m.SetMotorDuty(GetChannelMapLeftUp(), d.LeftUp)
	m.SetMotorDuty(GetChannelMapLeftLow(), d.LeftLow)
	m.SetMotorDuty(GetChannelMapRightUp(), d.RightUp)
	m.SetMotorDuty(GetChannelMapRightLow(), d.RightLow)
}
