package sacn

import (
	"fmt"
	"net"
	"v2sdl/config"
	"v2sdl/dmx"

	"github.com/rs/zerolog/log"
)

const (
	srvAddrTemplate = "239.255.%d.%d:5568"
	maxDatagramSize = 8192
)

// E131Packet defined by ANSI E1.31 2016 (c) ESTA - Section 4
type E131Packet struct {
	RootLayer
	FramingLayer
	DMPLayer
}

// RootLayer is from ANSI E1.31 2016 - Section 5
type RootLayer struct {
	PreambleSize      int16
	PostambleSize     int16
	APacketIdentifier [12]byte
	FlagsAndLength    int16
	Vector            int32
	CID               [16]byte
}

type FramingLayer struct {
	FlageAndLength int16
	Vector         int32
	SourceName     [32]byte
	Priority       byte
	SyncAddress    int16
	SeqenceNumber  byte
	Options        byte
	Universe       int16
}

type DMPLayer struct {
	FlagsAndLength int16
	Vector         byte
	AddressType    byte
	FirstProperty  int16
	AddressInc     int16
	PropertyCount  int16
	StartByte      byte
	Data           []byte
}

// SACN implements a NetDMX Listener
type SACN struct {
	dmx.Common
	socket *net.UDPConn
}

// NewService creates a new instance
func NewService() *SACN {
	return &SACN{}
}

func (x *SACN) Start(c *config.Config) (err error) {
	x.Cfg = c
	x.Universe = c.Universe
	x.Address = c.Address
	univHigh := c.Universe >> 8
	univLow := c.Universe & 255
	network := fmt.Sprintf(srvAddrTemplate, univHigh, univLow)

	ifi, err := net.InterfaceByName(c.Interface)
	if err != nil {
		log.Error().Err(err).Msg("Error")
		ifi = nil
	}
	gaddr, err := net.ResolveUDPAddr("udp", network)
	if err != nil {
		log.Error().Err(err).Msg("Error")
	}
	socket, err := net.ListenMulticastUDP("udp", ifi, gaddr)
	if err != nil {
		log.Error().Err(err).Msg("Error")
	}

	x.socket = socket

	if err == nil {
		go x.run()
	}
	return err
}

// Run starts a listening thread
func (x *SACN) run() {
	log.Info().Msg("Started E1.31 goroutine")
	defer log.Info().Msg("Exit goroutine")

	b := make([]byte, maxDatagramSize)
	for {
		n, addr, err := x.socket.ReadFrom(b)
		if err != nil {
			log.Error().Err(err).Msg("Error")
			break
		}

		if n < 0x7d { // Too small
			continue
		}

		// ETC Visualization Mode filter
		if b[0x7d] > 0 {
			continue
		}

		x.OnFrame(addr, b[0x7e:n])
	}
	x.socket.Close()
}

// Stop ends the running thread
func (x *SACN) Stop() {
	if x.socket != nil {
		x.socket.Close()
	}
}

func (x *SACN) Refresh(*config.Config) (err error) {
	return
}

func (x *SACN) Name() string {
	return fmt.Sprintf("sACN service")
}
