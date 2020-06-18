package ytcv

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type YTCV struct {
	client *http.Client
}

func New() *YTCV {

	return &YTCV{}
}

func getStringInBetween(str, start, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}

func (y *YTCV) FetchAll() ([]GridVideoRenderer, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}

	req, err := http.NewRequest("GET", "https://www.youtube.com/channel/UCCzUftO8KOVkV4wQG1vkUvg/videos", nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"user-agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	initialData := getStringInBetween(string(b), `window["ytInitialData"] = `, ";\n")
	var idr InitialDataResponse
	if err = json.Unmarshal([]byte(initialData), &idr); err != nil {
		return nil, err
	}

	var (
		continuation string
		itct         string
	)
	videos := make([]GridVideoRenderer, 0, 512)
	renderer := idr.Contents.TwoColumnBrowseResultsRenderer.Tabs[1].TabRenderer.Content.SectionListRenderer.Contents[0].ItemSectionRenderer.Contents[0].GridRenderer
	for _, item := range renderer.Items {
		videos = append(videos, item.GridVideoRenderer)
	}

	data := renderer.Continuations[0].NextContinuationData
	continuation = data.Continuation
	itct = data.ClickTrackingParams

	for continuation != "" {

		req, err = http.NewRequest("GET", "https://www.youtube.com/browse_ajax", nil)
		req.Header = http.Header{
			"x-youtube-client-name":    {"1"},
			"x-youtube-client-version": {"2.20200617.02.00"},
		}
		v := url.Values{
			"ctoken":       {continuation},
			"continuation": {continuation},
			"itct":         {itct},
		}
		req.URL.RawQuery = v.Encode()
		if err != nil {
			return nil, err
		}

		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var r ChannelVideosResponse
		if err = json.Unmarshal(b, &r); err != nil {
			return nil, err
		}

		gc := r[1].Response.ContinuationContents.GridContinuation
		for _, item := range gc.Items {
			videos = append(videos, item.GridVideoRenderer)
		}

		if len(gc.Continuations) < 1 {
			break
		}
		data = gc.Continuations[0].NextContinuationData
		continuation = data.Continuation
		itct = data.ClickTrackingParams
	}

	return videos, nil
}
