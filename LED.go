package main

// TODO: https://pkg.go.dev/gobot.io/x/gobot/drivers/gpio#RgbLedDriver
// Maybe this?

import (
	"log"

	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

func rgbToColor(r int, g int, b int) uint32 {
	return uint32(uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

type RGBStrip struct {
	ws2811 *ws281x.WS2811
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
		ws2811}

}
