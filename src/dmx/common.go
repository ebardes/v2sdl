package dmx

import (
	"bytes"
	"net"
	"v2sdl/config"
	"v2sdl/display"

	"github.com/rs/zerolog/log"
)

// NetDMX is the base interface for DMX
type NetDMX interface {
	config.Service
	SetDisplay(d *display.Display)
}

// Common stores aspects common to all Networked DMX implementations
type Common struct {
	NetDMX
	frame    []byte
	Cfg      *config.Config
	Universe int
	Address  int
	display  *display.Display
}

// OnFrame is the main event listener for when DMX packets arrive
func (me *Common) OnFrame(addr net.Addr, b []byte) {
	if !bytes.Equal(b, me.frame) {
		if len(me.frame) != len(b) {
			me.frame = make([]byte, len(b))
		}
		copy(me.frame, b)
		if me.Cfg.DebugLevel > 4 {
			log.Printf("Packet from %v size %d\n", addr, len(b))
		}

		if me.display != nil {
			me.display.OnFrame(b[me.Cfg.Address-1:])
		}
	}
}

func (me *Common) SetDisplay(d *display.Display) { me.display = d }
