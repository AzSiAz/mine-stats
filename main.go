package main

import (
	"log"
)

func main() {
	store, err := NewStore("db.bbolt")
	if err != nil {
		log.Fatalln("Error opening DB")
	}
	defer store.Close()
	log.Println("Done opening database")

	testMinecraft()
}

func testMinecraft() {
	mc := NewMinecraftServer("play.derentis.fr", 25565, 10)

	err := mc.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer mc.Socket.Close()

	log.Println("Connection to play.derentis.fr opened")
}
