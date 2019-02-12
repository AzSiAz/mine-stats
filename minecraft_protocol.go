package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

// Doc: https://wiki.vg/Query
// Go implementation: https://github.com/SpencerSharkey/gomc
// PHP implementation: https://github.com/xPaw/PHP-Minecraft-Query

type BasicMinecraftStatus struct {
}

type FullMinecraftStatus struct {
	BasicMinecraftStatus
}

type MinecraftServer struct {
	Address string
	Port    int64
	Timeout time.Duration
	Socket  net.Conn
}

func NewMinecraftServer(address string, port int64, timeout time.Duration) *MinecraftServer {
	return &MinecraftServer{
		Address: address,
		Port:    port,
		Timeout: timeout,
	}
}

func (sm *MinecraftServer) Connect() error {
	link := fmt.Sprintf("%v:%d", sm.Address, sm.Port)
	sock, err := net.Dial("udp", link)
	if err != nil {
		log.Println(err)
		return errors.New("error dialing udp server")
	}
	sm.Socket = sock

	return nil
}

func (sm *MinecraftServer) QueryBasic() (*BasicMinecraftStatus, error) {
	return &BasicMinecraftStatus{}, nil
}

func (sm *MinecraftServer) QueryFull() (*FullMinecraftStatus, error) {
	return &FullMinecraftStatus{}, nil
}
