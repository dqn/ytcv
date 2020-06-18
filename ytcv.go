package ytcv

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type YTCV struct {
	client *http.Client
}

func New() *YTCV {
	return &YTCV{}
}

func (y *YTCV) FetchAll() (interface{}, error) {
	req, err := http.NewRequest("GET", "", nil)
	v := url.Values{}
	req.URL.RawQuery = v.Encode()
	if err != nil {
		return nil, err
	}

	resp, err := y.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))
	return nil, nil
}
