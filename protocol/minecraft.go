package protocol

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mine-stats/models"
	"mine-stats/protocol/packet"
	"net"
	"time"
)

// Doc Query: https://wiki.vg/Query
// Doc Simple status ping: https://wiki.vg/Server_List_Ping
// Go implementation: https://github.com/SpencerSharkey/gomc
// PHP implementation: https://github.com/xPaw/PHP-Minecraft-Query

type MinecraftServer struct {
	Name string
	Address string
	Port    uint16
	Timeout time.Duration // in millisecond
}

func NewMinecraftServer(name string, address string, port uint16, timeout time.Duration) *MinecraftServer {
	return &MinecraftServer{
		Name: name,
		Address: address,
		Port:    port,
		Timeout: timeout,
	}
}

func (sm *MinecraftServer) Connect() (sock net.Conn, err error) {
	completeAddr := fmt.Sprintf("%s:%d", sm.Address, sm.Port)

	sock, err = net.Dial("tcp", completeAddr)
	if err != nil {
		return nil, errors.New("error dialing tcp server: " + err.Error())
	}

	return
}

func (sm *MinecraftServer) Close(sock net.Conn) error {
	err := sock.Close()
	return err
}

func (sm *MinecraftServer) SendHandshake(sock net.Conn) error {
	var handShakeBuffer []byte
	handShakeBuffer = append(handShakeBuffer, packet.PackVarInt(int32(47))...)
	handShakeBuffer = append(handShakeBuffer, packet.PackString(sm.Address)...)
	handShakeBuffer = append(handShakeBuffer, packet.PackUint16(sm.Port)...)
	handShakeBuffer = append(handShakeBuffer, byte(1))
	handshakePacket := &packet.Packet{
		ID:   0,
		Data: handShakeBuffer,
	}
	
	_, err := sock.Write(handshakePacket.Pack(-1))
	if err != nil {
		return errors.New("error writing handshake packet: " + err.Error())
	}

	return nil
}

func (sm *MinecraftServer) SendListPacket(sock net.Conn) error {
	listPacket := &packet.Packet{
		ID:   0,
		Data: []byte{},
	}
	_, err := sock.Write(listPacket.Pack(-1))
	if err != nil {
		return errors.New("error sending list packet")
	}

	return nil
}

func (sm *MinecraftServer) Query() (mc *models.MinecraftStatus, err error) {
	mc = &models.MinecraftStatus{
		Hostname: sm.Address,
		Port: sm.Port,
	}

	sock, err := sm.Connect()
	if err != nil {
		return
	}
	defer sm.Close(sock)

	err = sm.SendHandshake(sock)
	if err != nil { return }
	err = sm.SendListPacket(sock)
	if err != nil { return }

	defer sock.SetDeadline(time.Time{})
	err = sock.SetDeadline(time.Now().Add(sm.Timeout * time.Second))
	if err != nil { return }

	received, err := packet.RecvPacket(bufio.NewReader(sock), false)
	if err != nil {
		return nil, errors.New("failed to received list packet: " + err.Error())
	}
	s, err := packet.UnpackString(bytes.NewReader(received.Data))
	if err != nil {
		return nil, errors.New("Error unpacking data" + err.Error())
	}

	err = json.Unmarshal([]byte(s), mc)
	if err != nil {
		return nil, errors.New("error unmarshalling json from server: " + err.Error())
	}

	return
}