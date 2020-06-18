# ytcv

Fetch YouTube videos from a specific channel without authentication.

## Installation

```bash
$ go get github.com/dqn/ytcv
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/dqn/ytcv"
)

func main() {
	videos, err := ytcv.FetchAll("CHANNEL_ID")
	if err != nil {
		// handle error
	}

	for _, video := range videos {
		fmt.Println(video.VideoID)
		fmt.Println(video.Title.SimpleText)
		fmt.Println(video.PublishedTimeText.SimpleText)
		fmt.Println(video.Thumbnail.Thumbnails[0].URL)
		fmt.Println(video.ViewCountText.SimpleText)
	}
}
```

## License

MIT
