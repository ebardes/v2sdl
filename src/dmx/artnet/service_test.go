package artnet

import (
	"testing"
	"v2sdl/config"
)

func TestNode(t *testing.T) {
	cfg := config.Config{}
	s, err := NewService(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	s.run()
}
