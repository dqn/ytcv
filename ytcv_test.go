package ytcv

import "testing"

func TestFetchALl(t *testing.T) {
	y := New()
	_, err := y.FetchAll()
	if err != nil {
		t.Fatal(err)
	}
}
