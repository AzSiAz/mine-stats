package minecraftProtocol

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"math"
)

// Packet define a net data package
type Packet struct {
	ID   byte
	Data []byte
}

// Pack
func (p *Packet) Pack(threshold int) (pack []byte) {
	data := []byte{p.ID}
	data = append(data, p.Data...)

	if threshold > 0 {
		if len(data) > threshold {
			Len := len(data)
			VarLen := PackVarInt(int32(Len))
			data = Compress(data)

			pack = append(pack, PackVarInt(int32(len(VarLen)+len(data)))...)
			pack = append(pack, VarLen...)
			pack = append(pack, data...)
		} else {
			pack = append(pack, PackVarInt(int32(len(data)+1))...)
			pack = append(pack, 0x00)
			pack = append(pack, data...)
		}
	} else {
		pack = append(pack, PackVarInt(int32(len(data)))...) //len
		pack = append(pack, data...)
	}

	return
}

// PackString
func PackString(s string) (p []byte) {
	byteString := []byte(s)
	p = append(p, PackVarInt(int32(len(byteString)))...) //len
	p = append(p, byteString...)                         //data
	return
}

// PackUint16
func PackUint16(n uint16) (p []byte) {
	p = []byte{
		byte(n >> 8),
		byte(n),
	}
	return
}

// PackUint32
func PackUint32(n uint32) (p []byte) {
	p = []byte{
		byte(n >> 24), byte(n >> 16),
		byte(n >> 8), byte(n),
	}
	return
}

// PackUint64
func PackUint64(n uint64) (p []byte) {
	p = []byte{
		byte(n >> 56), byte(n >> 48), byte(n >> 40), byte(n >> 32),
		byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n),
	}
	return
}

//PackVarInt
func PackVarInt(n int32) (VarInt []byte) {
	num := uint32(n)
	for {
		b := num & 0x7F
		num >>= 7
		if num != 0 {
			b |= 0x80
		}
		VarInt = append(VarInt, byte(b))
		if num == 0 {
			break
		}
	}
	return
}

//PackPosition
func PackPosition(x, y, z int) (p []byte) {
	p = make([]byte, 8)
	position := (int64(x&0x3FFFFFF) << 38) | int64((y&0xFFF)<<26) | int64(z&0x3FFFFFF)
	for i := 7; i >= 0; i-- {
		p[i] = byte(position)
		position >>= 8
	}
	return
}

//PackFloat
func PackFloat(f float32) (p []byte) {
	return PackUint32(math.Float32bits(f))
}

//PackDouble
func PackDouble(d float64) (p []byte) {
	return PackUint64(math.Float64bits(d))
}

//PackBoolean
func PackBoolean(b bool) byte {
	if b {
		return 0x01
	}
	return 0x00
}

//ReadNBytes read N bytes from bytes.Reader
func ReadNBytes(b *bytes.Reader, n int) (bs []byte, err error) {
	bs = make([]byte, n)
	for i := 0; i < n; i++ {
		bs[i], err = b.ReadByte()
		if err != nil {
			return
		}
	}
	return
}

//UnpackString
func UnpackString(b *bytes.Reader) (s string, err error) {
	l, err := UnpackVarInt(b)
	if err != nil {
		return "", err
	}

	bs, err := ReadNBytes(b, int(l))

	return string(bs), err
}

//UnpackVarInt
func UnpackVarInt(b *bytes.Reader) (int32, error) {
	var n uint
	for i := 0; i < 5; i++ {
		sec, err := b.ReadByte()
		if err != nil {
			return 0, err
		}

		n |= (uint(sec&0x7F) << uint(7*i))

		if sec&0x80 == 0 {
			break
		}
	}
	return int32(n), nil
}

//UnpackInt16
func UnpackInt16(b *bytes.Reader) (int16, error) {
	bs, err := ReadNBytes(b, 2)
	if err != nil {
		return 0, err
	}
	return int16(bs[0])<<8 | int16(bs[1]), nil
}

//UnpackInt32
func UnpackInt32(b *bytes.Reader) (int32, error) {
	bs, err := ReadNBytes(b, 4)
	if err != nil {
		return 0, err
	}
	return int32(bs[0])<<24 | int32(bs[1])<<16 | int32(bs[2])<<8 | int32(bs[3]), nil
}

//UnpackInt64
func UnpackInt64(b *bytes.Reader) (int64, error) {
	bs, err := ReadNBytes(b, 8)
	if err != nil {
		return 0, err
	}
	return int64(bs[0])<<56 | int64(bs[1])<<48 | int64(bs[2])<<40 | int64(bs[3])<<32 |
		int64(bs[4])<<24 | int64(bs[5])<<16 | int64(bs[6])<<8 | int64(bs[7]), nil
}

// UnpackPosition
func UnpackPosition(b *bytes.Reader) (x, y, z int, err error) {
	position, err := UnpackInt64(b)

	x = int(position >> 38)
	y = int((position >> 26) & 0xFFF)
	z = int(position << 38 >> 38)

	if x >= 1<<25 {
		x -= 1 << 26
	}
	if y >= 1<<11 {
		y -= 1 << 12
	}
	if z >= 1<<25 {
		z -= 1 << 26
	}
	return
}

// UnpackFloat
func UnpackFloat(b *bytes.Reader) (float32, error) {
	n, err := UnpackInt32(b)
	return math.Float32frombits(uint32(n)), err
}

// UnpackDouble
func UnpackDouble(b *bytes.Reader) (float64, error) {
	n, err := UnpackInt64(b)
	return math.Float64frombits(uint64(n)), err
}

// RecvPacket recive a packet from server
func RecvPacket(s *bufio.Reader, useZlib bool) (*Packet, error) {
	var len int
	for i := 0; i < 5; i++ {
		b, err := s.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read len of packet fail: %v", err)
		}
		len |= (int(b&0x7F) << uint(7*i))
		if b&0x80 == 0 {
			break
		}
	}

	if len < 1 {
		return nil, fmt.Errorf("packet length too short")
	}

	data := make([]byte, len)
	var err error
	for i := 0; i < len; i++ {
		data[i], err = s.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read content of packet fail: %v", err)
		}
	}

	if useZlib {
		return UnCompress(data)
	}

	return &Packet{
		ID:   data[0],
		Data: data[1:],
	}, nil
}

func UnCompress(data []byte) (*Packet, error) {
	reader := bytes.NewReader(data)
	sizeUncompressed, err := UnpackVarInt(reader)
	if err != nil {
		return nil, err
	}

	uncompressData := make([]byte, sizeUncompressed)
	if sizeUncompressed != 0 { // != 0 means compressed, let's decompress
		r, err := zlib.NewReader(reader)

		if err != nil {
			return nil, fmt.Errorf("decompress fail: %v", err)
		}
		_, err = io.ReadFull(r, uncompressData)
		if err != nil {
			return nil, fmt.Errorf("decompress fail: %v", err)
		}
		r.Close()
	} else {
		uncompressData = data[1:]
	}
	return &Packet{
		ID:   uncompressData[0],
		Data: uncompressData[1:],
	}, nil
}

func Compress(data []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(data) // nolint;
	w.Close()
	return b.Bytes()
}
