package video

import "testing"

func Test1(t *testing.T) {
	fn := "/Users/eric/Downloads/doctor.who.2005.s11e00.twice.upon.a.time.christmas.special.1080p.hdtv.h264-mtb.mkv"
	a, _ := NewAVMediaPlayer()
	defer a.Dispose()

	err := a.Open(fn)
	if err != nil {
		t.Error(err)
	}
}
