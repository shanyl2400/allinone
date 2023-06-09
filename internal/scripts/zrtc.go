package scripts

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	zrtcURL = "http://192.168.0.210/zrtc/"
	re      = regexp.MustCompile(`<a href=\"(zrtc.*)\">(zrtc.*)</a>`)

	ErrNoSuchZRTC = errors.New("no such zrtc")
)

type ZRTC struct {
	Name string
	Path string
}

type ZrtcOp struct {
	client http.Client
}

func (z *ZrtcOp) ListZrtcs() ([]*ZRTC, error) {
	req, err := http.NewRequest("GET", zrtcURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := z.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	zrtcList := make([]*ZRTC, 0)
	match := re.FindAllStringSubmatch(string(bytes), -1)

	for _, d := range match {
		zrtcItem := &ZRTC{
			Name: d[1],
			Path: d[2],
		}
		zrtcList = append(zrtcList, zrtcItem)
	}
	return zrtcList, nil
}

func (z *ZrtcOp) DownloadZrtc(path string) ([]string, error) {
	zrtcList, err := z.ListZrtcs()
	if err != nil {
		return nil, err
	}
	ans := make([]string, 0)

	flag := false
	for _, item := range zrtcList {
		if item.Path == path {
			flag = true
			break
		}
	}
	if !flag {
		log.Printf("No such zrtc: %v", path)
		return nil, ErrNoSuchZRTC
	}

	//对.tar.gz包解压
	if strings.HasSuffix(path, ".tar.gz") {
		output, err := execute("wget", zrtcURL+path, "-O", "zrtc.tar.gz", "-q")
		ans = append(ans, output)
		if err != nil {
			return ans, err
		}

		ans = append(ans, "tar xvf zrtc.tar.gz")
		output, err = execute("tar", "xvf", "zrtc.tar.gz")
		ans = append(ans, output)
		if err != nil {
			return ans, err
		}

		ans = append(ans, "rm -f zrtc.tar.gz")
		output, err = execute("rm", "-f", "zrtc.tar.gz")
		ans = append(ans, output)
		if err != nil {
			return ans, err
		}
	} else {
		output, err := execute("wget", zrtcURL+path, "-O", "zrtc", "-q")
		ans = append(ans, output)
		if err != nil {
			return ans, err
		}
	}

	output2, err := execute("chmod", "+x", "zrtc")
	ans = append(ans, output2)
	if err != nil {
		return ans, err
	}
	return ans, nil
}
