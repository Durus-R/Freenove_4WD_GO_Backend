package car

import (
	"errors"
	"github.com/go-daq/smbus"
	"log"
	"math"
)

type ADC struct {
	address    uint8
	pcf8591Cmd uint8
	ads7830Cmd uint8
	index      string
	bus        *smbus.Conn
}

func NewADC() *ADC {
	conn, err := smbus.Open(1, 0x48)

	defer func(conn *smbus.Conn) {
		err := conn.Close()
		if err != nil {
			log.Print("Error at closing: ", err)
		}
	}(conn)

	if err != nil {
		log.Fatal("Cannot open bus. "+
			"Make sure it is properly connected. "+
			"\nError output: ", err)
	}
	var index string
	for i := 0; i < 3; i++ {
		aa, err := conn.ReadReg(0x48, 0xf4)
		if err != nil {
			log.Print("Error: ", err)
		}
		if aa < 150 {
			index = "PCF8591"
		} else {
			index = "ADS7830"
		}
	}
	return &ADC{
		address:    0x48,
		pcf8591Cmd: 0x40,
		ads7830Cmd: 0x84,
		index:      index,
		bus:        conn,
	}

}

func sortArray(arr [9]uint8) [9]uint8 {
	for i := 0; i <= len(arr)-1; i++ {
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

func (a *ADC) analogReadPCF8591(channel uint8) (uint8, error) {
	measures := [9]uint8{}
	err := errors.New("") // Empty Error
	for i := 0; i < 9; i++ {
		measures[i], err = a.bus.ReadReg(a.address,
			a.pcf8591Cmd+channel)
		if err != nil {
			log.Print("Error while reading: ", err)
			return 0, err
		}
	}
	sorted := sortArray(measures)
	return sorted[4], nil
}

func (a *ADC) receivePCF8591(channel uint8) (float64, error) {
	var value1 uint8
	var value2 uint8
	for {
		value1, _ = a.analogReadPCF8591(channel)
		value2, _ = a.analogReadPCF8591(channel)
		if value1 == value2 {
			if value1 == 0 {
				return 0, errors.New("something failed while reading")
			}
			break
		}
	}
	result := float64(value1) / 256.0 * 3.3
	return math.Round(result*100) / 100, nil
}

func (a *ADC) receiveADS7830(channel uint8) (float64, error) {
	commandSet := a.ads7830Cmd | ((((channel << 2) | (channel >> 1)) & 0x07) << 4)
	_, err := a.bus.WriteByte(commandSet)
	if err != nil {
		log.Print("Error writing to bus: ", err)
		return 0, err
	}
	var value1 uint8
	var value2 uint8
	for {
		value1, _ = a.bus.ReadReg(a.address, 0)
		value2, _ = a.bus.ReadReg(a.address, 0)
		if value1 == value2 {
			break
		}
	}
	result := float64(value1) / 255 * 3.3
	return math.Round(result*100) / 100, nil
}

func (a *ADC) receiveADC(channel uint8) (float64, error) {
	var data float64
	var err error
	switch a.index {
	case "PCF8591":
		data, err = a.receivePCF8591(channel)
	case "ADS7830":
		data, err = a.receiveADS7830(channel)
	}
	return data, err
}

func (a *ADC) Battery() (float64, error) {
	res, err := a.receiveADC(2)
	return res * 8, err
}

func (a *ADC) IDR() ([2]float64, error) {
	var (
		left  float64
		right float64
		err   error
	)
	left, err = a.receiveADC(0)
	if err != nil {
		return [2]float64{}, err
	}
	right, err = a.receiveADC(1)
	if err != nil {
		return [2]float64{}, err
	}
	return [2]float64{
		left,
		right,
	}, nil
}
