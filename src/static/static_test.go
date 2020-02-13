package static

import "testing"

func TestLink(t *testing.T) {
	item, err := GetItem("index.html")

	if err != nil {

	}
	t.Log(item, err)

	item, err = GetItem("doesn't exist at all")
	t.Log(item, err)
}
