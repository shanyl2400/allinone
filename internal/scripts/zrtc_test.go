package scripts

import "testing"

func TestListZRTC(t *testing.T) {
	c := new(ZrtcOp)
	zrtcs, err := c.ListZrtcs()
	if err != nil {
		t.Error(err)
	}
	for i := range zrtcs {
		t.Log(zrtcs[i].Name, zrtcs[i].Path)
	}
}

func TestDownloadZRTC(t *testing.T) {
	c := new(ZrtcOp)
	zrtcs, err := c.ListZrtcs()
	if err != nil {
		t.Error(err)
	}

	path := ""
	for i := range zrtcs {
		if zrtcs[i].Path == "zrtc23.3.9.24p" {
			path = zrtcs[i].Path
		}
	}

	output, err := c.DownloadZrtc(path)
	if err != nil {
		t.Error(err)
	}
	t.Log(output)
}
