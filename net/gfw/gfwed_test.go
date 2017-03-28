package gfw

import (
	"testing"
)

func TestIsGFWed(t *testing.T) {

	gfwlist, err := CreateGFWList(gfwlistURL, "gfwlist.lst")
	if err != nil {
		t.Error(err)
	}

	if gfwlist.Hit("twitter.com") == false {
		t.Error("twitter.com should be gfwed")
	}
}
