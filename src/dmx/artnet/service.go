package artnet

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"v2sdl/config"
	"v2sdl/dmx"
)

const (
	ArtNetPort = 6454

	OpPoll       = 0x2000
	OpPollReply  = 0x2100
	OpDmx        = 0x5000
	OpNzs        = 0x5100
	OpTodRequest = 0x8000
	OpTodData    = 0x8100
	OpRdm        = 0x8300
	OpRdmSub     = 0x8400
)

type ArtNet struct {
	dmx.Common
	Sock *net.UDPConn
}

type ArtNetPacket struct {
	Header [8]byte
	OpCode int16
}

type ArtPollReplyPacket struct {
	ArtNetPacket
	IP          [4]byte
	Port        int16
	VersInfoH   byte
	VersInfoL   byte
	NetSwitch   byte
	SubSwitch   byte
	OemHi       byte
	Oem         byte
	UbeaVersion byte
	Status1     byte
	EstaManLo   byte
	EstaManHi   byte
	ShortName   [18]byte
	LongName    [64]byte
	NodeReport  [64]byte
	NumPortsHi  byte
	NumPortsLo  byte
	PortTypes   [4]byte
	GoodInput   [4]byte
	GoodOutput  [4]byte
	SwIn        [4]byte
	SwOut       [4]byte
	SwVideo     byte
	SwMacro     byte
	SwRemote    byte
	Spare       [3]byte
	Style       byte
	MAC         [6]byte
	BindIP      [4]byte
	BindIndex   byte
	Status2     byte
	Filler      [26]byte
}

type ArtDmxPacket struct {
	ArtNetPacket
	ProtVerHi byte
	ProtVerLo byte
	Sequence  byte
	Physical  byte
	SubUni    byte
	Net       byte
	LengthHi  byte
	LengthLo  byte
	StartByte byte
	Data      [512]byte
}

type ArtRdm struct {
	ArtNetPacket
	ProtVerHi byte
	ProtVerLo byte
	RdmVer    byte
	Filler2   byte
	Spare     [7]byte
	Net       byte
	Command   byte
	Address   byte
	RdmPacket []byte
}

type ArtRdmSub struct {
	ArtNetPacket
	ProtVerHi    byte
	ProtVerLo    byte
	RdmVer       byte
	Filler2      byte
	UID          [6]byte
	Spare1       byte
	CommandClass byte
	ParameterId  int16
	SubDevice    int16
	SubCount     int16
	Spare2       byte
	Spare3       byte
	Spare4       byte
	Spare5       byte
	Data         []byte
}

type ArtTodRequest struct {
	ArtNetPacket
	ProtVerHi byte
	ProtVerLo byte
	Filler1   byte
	Filler2   byte
	Spare     [7]byte
	Net       byte
	Command   byte
	AddCount  byte
}

// decode streams binary data into structures
func decode(b []byte, data interface{}) (r io.Reader, err error) {
	r = bytes.NewBuffer(b)
	err = binary.Read(r, binary.LittleEndian, data)
	return
}

func u16(b []byte, offset int) uint {
	return uint(b[offset]) + (uint(b[offset+1]) << 8)
}

func NewService(c *config.Config) (*ArtNet, error) {
	var err error

	laddr := &net.UDPAddr{Port: ArtNetPort}
	ifi, err := net.InterfaceByName(c.Interface)
	if err != nil {
		log.Println(err)
		ifi = nil
	} else {
		addr, err := ifi.Addrs()
		if err != nil {
			log.Println(err)
		}
		for _, ip := range addr {
			log.Println(ip.Network())
		}
	}
	a := &ArtNet{}
	a.Cfg = c
	a.Sock, err = net.ListenUDP("udp", laddr)
	return a, err
}

func (x *ArtNet) Start(*config.Config) error {
	go x.run()
	return nil
}

func (x *ArtNet) run() {
	log.Println(x)
	b := make([]byte, 8192)
	for {
		i, addr, err := x.Sock.ReadFrom(b)
		if err != nil {
			log.Fatal(err)
		}

		opCode := u16(b, 8)

		switch opCode {
		case OpDmx:
			var d ArtDmxPacket
			decode(b, &d)

			// ETC Visualization Mode filter
			if d.Data[0] > 0 {
				continue
			}

			x.OnFrame(addr, d.Data[:])
		case OpPoll:
			x.Poll(b)

		case OpTodRequest:
			x.DoOpTodRequest(b[:i])

		case OpPollReply:
		default:
			fmt.Printf("OpCode=%04x Sock=%v\n", opCode, addr)
			fmt.Println(hex.Dump(b[:i]))
		}
	}
}

func (x *ArtNet) Stop() {

}

func (x *ArtNet) Refresh(*config.Config) (err error) {
	return
}

func (x *ArtNet) DoOpTodRequest(b []byte) {
	var d ArtTodRequest

	r, err := decode(b, &d)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Hi:%d, Lo:%d add Net:%d\n", d.ProtVerHi, d.ProtVerLo, d.Net)
		fmt.Printf("%v\n", d)

		var b1 [1]byte
		for i := byte(0); i < d.AddCount; i++ {
			fmt.Println(r.Read(b1[:]))
		}
	}
}

func (x *ArtNet) PollReply() {
	p := ArtPollReplyPacket{}

	copy(p.Header[0:], []byte{65, 114, 116, 45, 78, 101, 116, 0})
	p.OpCode = 0x2100
	p.Port = ArtNetPort

	copy(p.LongName[0:], []byte("Bardes Media"))
	copy(p.ShortName[0:], []byte("Bardes Media"))
	copy(p.NodeReport[0:], []byte("Copacetic"))

	b := make([]byte, 0)
	bb := bytes.NewBuffer(b)

	if err := binary.Write(bb, binary.LittleEndian, p); err != nil {
		fmt.Println(err)
	}

	// fmt.Println(hex.Dump(bb.Bytes()))
	x.Sock.Write(bb.Bytes())
}

func (x *ArtNet) Poll(b []byte) {
	var p ArtPollReplyPacket

	if err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &p); err != nil {
		fmt.Println(err)
	} else {
		x.PollReply()
	}
}
