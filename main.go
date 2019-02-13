package main

import (
	"github.com/Tnze/gomcbot"
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
	test2()
}

func test2() {
	resp, err := gomcbot.PingAndList("mc.hypixel.net", 25565)
	if err != nil {
		log.Fatalln(err)
	}

	println(resp)
}

func testMinecraft() {
	mc := NewMinecraftServer("mc.hypixel.net", 25565, 30)

	err := mc.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer mc.Socket.Close()
	log.Println("Connection to mc.hypixel.net opened")

	status, err := mc.QueryBasic()
	if err != nil {
		log.Fatalln(err)
	}

	println("Hostname: " + status.Hostname)
}
