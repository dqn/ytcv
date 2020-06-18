package ytcv

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
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

func (y *YTCV) FetchAll() (interface{}, error) {
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

	for _, video := range videos {
		println(video.VideoID)
	}
	println(len(videos))
	println(continuation)
	println(itct)

	// req, err := http.NewRequest("GET", "https://www.youtube.com/browse_ajax", nil)
	// req.Header = http.Header{
	// 	"cookie":                      {"YSC=6defHPC6P_4; GPS=1; VISITOR_INFO1_LIVE=jfPtGxbPI9U"},
	// 	"referer":                     {"https://www.youtube.com/channel/UCCzUftO8KOVkV4wQG1vkUvg/videos"},
	// 	"sec-fetch-dest":              {"empty"},
	// 	"sec-fetch-mode":              {"cors"},
	// 	"sec-fetch-site":              {"same-origin"},
	// 	"user-agent":                  {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"},
	// 	"x-spf-previous":              {"https://www.youtube.com/channel/UCCzUftO8KOVkV4wQG1vkUvg/videos"},
	// 	"x-spf-referer":               {"https://www.youtube.com/channel/UCCzUftO8KOVkV4wQG1vkUvg/videos"},
	// 	"x-youtube-ad-signals":        {"dt=1592487601145&flash=0&frm&u_tz=540&u_his=2&u_java&u_h=1080&u_w=1920&u_ah=1010&u_aw=1920&u_cd=24&u_nplug=3&u_nmime=4&bc=31&bih=899&biw=1388&brdim=0%2C23%2C0%2C23%2C1920%2C23%2C1920%2C1010%2C1403%2C899&vis=1&wgl=true&ca_type=image&bid=ANyPxKpm8828vEFaA3gk7n3HRI76ki9xeX6-OXIw2cIcEF21NKDbZNThgHE5mWvaM9kM1_Bqj2dLkcOg9JUsvUnd85mRnCWn2Q"},
	// 	"x-youtube-client-name":       {"1"},
	// 	"x-youtube-client-version":    {"2.20200617.02.00"},
	// 	"x-youtube-device":            {"cbr=Chrome&cbrver=83.0.4103.97&ceng=WebKit&cengver=537.36&cos=Macintosh&cosver=10_15_5"},
	// 	"x-youtube-page-cl":           {"316682464"},
	// 	"x-youtube-page-label":        {"youtube.ytfe.desktop_20200616_2_RC0"},
	// 	"x-youtube-time-zone":         {"Asia/Tokyo"},
	// 	"x-youtube-utc-offset":        {"540"},
	// 	"x-youtube-variants-checksum": {"01793b9d32911a0f92629a032e992180"},
	// }
	// v := url.Values{
	// 	"ctoken":       {"4qmFsgI0EhhVQ0N6VWZ0TzhLT1ZrVjR3UUcxdmtVdmcaGEVnWjJhV1JsYjNNZ0FEZ0JlZ0V5dUFFQQ%3D%3D"},
	// 	"continuation": {"4qmFsgI0EhhVQ0N6VWZ0TzhLT1ZrVjR3UUcxdmtVdmcaGEVnWjJhV1JsYjNNZ0FEZ0JlZ0V5dUFFQQ%3D%3D"},
	// 	"itct":         {"CCsQybcCIhMI9YbNzb6L6gIVAQkqCh3uMg9l"},
	// }
	// req.URL.RawQuery = v.Encode()
	// if err != nil {
	// 	return nil, err
	// }

	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()

	// b, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Println(string(b))
	return nil, nil
}
