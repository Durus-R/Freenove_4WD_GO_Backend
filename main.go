package main

import (
	"log"
	"os/user"
)

// TODO: Ultrasonic, Camera
func main() {
	usr, _ := user.Current()
	if usr.Name != "root" {
		log.Fatal("Please restart this system with root access")
	}
}
