package main

import (
	"errors"
	"github.com/go-daq/smbus"
	"log"
	"math"
)

type ADC struct {
	Address    uint8
	Pcf8591Cmd uint8
	Ads7830Cmd uint8
	Index      string
	Bus        *smbus.Conn
}

func CreateADC() ADC {
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
	return ADC{
		Address:    0x48,
		Pcf8591Cmd: 0x40,
		Ads7830Cmd: 0x84,
		Index:      index,
		Bus:        conn,
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

func (a ADC) analogReadPCF8591(channel uint8) uint8 {
	measures := [9]uint8{}
	err := errors.New("") // Empty Error
	for i := 0; i < 9; i++ {
		measures[i], err = a.Bus.ReadReg(a.Address,
			a.Pcf8591Cmd+channel)
		if err != nil {
			log.Print("Error while reading: ", err)
		}
	}
	sorted := sortArray(measures)
	return sorted[4]
}

func (a ADC) receivePCF8591(channel uint8) float64 {
	var value1 uint8
	var value2 uint8
	for {
		value1 = a.analogReadPCF8591(channel)
		value2 = a.analogReadPCF8591(channel)
		if value1 == value2 {
			break
		}
	}
	result := float64(value1) / 256.0 * 3.3
	return math.Round(result*100) / 100
}

func (a ADC) receiveADS7830(channel uint8) float64 {
	commandSet := a.Ads7830Cmd | ((((channel << 2) | (channel >> 1)) & 0x07) << 4)
	_, err := a.Bus.WriteByte(commandSet)
	if err != nil {
		log.Print("Error writing to Bus: ", err)
	}
	var value1 uint8
	var value2 uint8
	for {
		value1, _ = a.Bus.ReadReg(a.Address, 0)
		value2, _ = a.Bus.ReadReg(a.Address, 0)
		if value1 == value2 {
			break
		}
	}
	result := float64(value1) / 255 * 3.3
	return math.Round(result*100) / 100
}

func (a ADC) receiveADC(channel uint8) float64 {
	var data float64
	switch a.Index {
	case "PCF8591":
		data = a.receivePCF8591(channel)
	case "ADS7830":
		data = a.receiveADS7830(channel)
	}
	return data
}

func (a ADC) Battery() float64 {
	return a.receiveADC(2) * 8
}

func (a ADC) IDR() [2]float64 {
	return [2]float64{
		a.receiveADC(0),
		a.receiveADC(1),
	}
}
