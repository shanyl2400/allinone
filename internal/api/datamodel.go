package api

type errResponse struct {
	Message string `json:"message"`
}

type listBranchResponse struct {
	Branches []string `json:"branches"`
}

type zrtc struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type publishRecord struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	GomssBranch string `json:"gomss_branch"`
	ZrtcVersion string `json:"zrtc_version"`
	RecordedAt  int64  `json:"recorded_at"`
}

type listZRTCResopnse struct {
	Zrtcs []*zrtc `json:"zrtcs"`
}

type listPublishRecordersResponse struct {
	Records []*publishRecord `json:"records"`
}

type getPublishLogsResonse struct {
	Logs []string `json:"logs"`
}

type publishRequest struct {
	GomssBranch string `json:"gomss_branch"`
	ZRTCPath    string `json:"zrtc_path"`
	Version     string `json:"version"`
	LocalZRTC   bool   `json:"local_zrtc"`
}
