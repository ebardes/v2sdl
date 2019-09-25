package artnet

import (
	"testing"
)

func TestNode(t *testing.T) {
	s, err := NewService(nil)
	if err != nil {
		t.Fatal(err)
	}
	s.Run()
	s.Stop()
}
