package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

// Doc Query: https://wiki.vg/Query
// Doc Simple status ping: https://wiki.vg/Server_List_Ping
// Go implementation: https://github.com/SpencerSharkey/gomc
// PHP implementation: https://github.com/xPaw/PHP-Minecraft-Query

type PlayerInfo struct {
	Max     int64
	Current int64
	Players []string
}

type BasicMinecraftStatus struct {
	Hostname      string
	GameType      string
	Map           string
	ServerVersion string
	Motd          string
	PlayerInfo    PlayerInfo
}

type FullMinecraftStatus struct {
	BasicMinecraftStatus
}

type MinecraftServer struct {
	Address string
	Port    uint16
	Timeout time.Duration // in millisecond
	Socket  net.Conn
}

func NewMinecraftServer(address string, port uint16, timeout time.Duration) *MinecraftServer {
	return &MinecraftServer{
		Address: address,
		Port:    port,
		Timeout: timeout,
	}
}

func (sm *MinecraftServer) Connect() error {
	completeAddr := fmt.Sprintf("%s:%d", sm.Address, sm.Port)

	sock, err := net.Dial("tcp", completeAddr)
	if err != nil {
		log.Println(err)
		return errors.New("error dialing udp server: " + err.Error())
	}
	sm.Socket = sock

	return nil
}

func (sm *MinecraftServer) QueryBasic() (*BasicMinecraftStatus, error) {
	mc := &BasicMinecraftStatus{}

	var handShakeBuffer []byte
	handShakeBuffer = append(handShakeBuffer, PackVarInt(int32(47))...)
	handShakeBuffer = append(handShakeBuffer, PackString(sm.Address)...)
	handShakeBuffer = append(handShakeBuffer, PackUint16(sm.Port)...)
	handShakeBuffer = append(handShakeBuffer, byte(1))
	handshakePacket := &Packet{
		ID:   0,
		Data: handShakeBuffer,
	}
	_, err := sm.Socket.Write(handshakePacket.Pack(-1))
	if err != nil {
		return nil, errors.New("error writing handshake packet")
	}

	listPacket := &Packet{
		ID:   0,
		Data: []byte{1},
	}
	_, err = sm.Socket.Write(listPacket.Pack(-1))
	if err != nil {
		return nil, errors.New("error sending list packet")
	}

	//defer sm.Socket.SetDeadline(time.Time{})
	sm.Socket.SetDeadline(time.Now().Add(sm.Timeout * time.Second))

	received, err := RecvPacket(bufio.NewReader(sm.Socket), false)
	if err != nil {
		return nil, errors.New("failed to received list packet: " + err.Error())
	}
	s, err := UnpackString(bytes.NewReader(received.Data))
	if err != nil {
		return nil, errors.New("Error unpacking data" + err.Error())
	}

	println(s)

	return mc, nil
}

func (sm *MinecraftServer) QueryFull() (*FullMinecraftStatus, error) {
	return &FullMinecraftStatus{}, nil
}

func scanDelimitedResponse(input []byte, eof bool) (adv int, token []byte, err error) {
	if len(input) == 0 {
		return 0, nil, errors.New("end of input")
	}
	i := bytes.Index(input, []byte{0x00})
	return i + 1, input[:i], nil
}
