package video

import "testing"

func Test1(t *testing.T) {
	fn := "/Users/eric/toppcell-code/images/oceans-clip.mp4"
	a, _ := NewAVMediaPlayer()
	defer a.Dispose()

	err := a.Open(fn)
	if err != nil {
		t.Error(err)
	}

	err = a.Run()
	if err != nil {
		t.Error(err)
	}
}
