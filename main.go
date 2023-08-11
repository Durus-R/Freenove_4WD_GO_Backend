package main

import (
	"log"
	"os/user"
)

func main() {
	usr, _ := user.Current()
	if usr.Name != "root" {
		log.Fatal("Please restart this system with root access")
	}
}
