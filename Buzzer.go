package main

import (
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
	"log"
	"strconv"
)

type Buzzer struct {
	Pin int
}

func SetupBuzzer() *gpio.BuzzerDriver {
	rpi := raspi.NewAdaptor()
	buzzer := gpio.NewBuzzerDriver(rpi, strconv.Itoa(17))
	err := buzzer.Start()
	if err != nil {
		log.Fatal("Error while Buzzing: ", err)
	}
	return buzzer
}
