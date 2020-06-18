package ytcv

import (
	"os"
	"testing"
)

func TestFetchAll(t *testing.T) {
	videos, err := FetchAll(os.Getenv("CHANNEL_ID"))
	if err != nil {
		t.Fatal(err)
	}

	for _, video := range videos {
		println(video.VideoID)
	}
	println(len(videos))
}
