package scripts

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var (
	zrtcURL = "https://go-test001.livecourse.com/zrtc/"
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

	// 这里要格式化再输出，因为 ReadAll 返回的是字节切片
	// fmt.Printf("%s\n", bytes)

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

func (z *ZrtcOp) DownloadZrtc(path string) (map[string]string, error) {
	zrtcList, err := z.ListZrtcs()
	if err != nil {
		return nil, err
	}
	ans := make(map[string]string)

	flag := false
	for _, item := range zrtcList {
		if item.Path == path {
			flag = true
			break
		}
	}
	if !flag {
		log.Fatalf("No such zrtc: %v", path)
		return nil, ErrNoSuchZRTC
	}

	output, err := execute("wget", zrtcURL+path, "-O", "zrtc0", "-q")
	ans["wget"] = output
	if err != nil {
		return ans, err
	}

	output2, err := execute("chmod", "+x", "zrtc0")
	ans["chmod"] = output2
	if err != nil {
		return ans, err
	}
	return ans, nil
}
