package ytcv

import "testing"

func TestFetchALl(t *testing.T) {
	y := New()
	videos, err := y.FetchAll()
	if err != nil {
		t.Fatal(err)
	}

	for _, video := range videos {
		println(video.VideoID)
	}
	println(len(videos))
}
