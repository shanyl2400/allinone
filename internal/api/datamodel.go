package api

type errResponse struct {
	Message string `json:"message"`
}

type listBranchResponse struct {
	Branches []string `json:"branches"`
}

type ZRTC struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type listZRTCResopnse struct {
	Zrtcs []*ZRTC `json:"zrtcs"`
}

type publishRequest struct {
	GomssBranch string `json:"gomss_branch"`
	ZRTCPath    string `json:"zrtc_path"`
	Version     string `json:"version"`
}
