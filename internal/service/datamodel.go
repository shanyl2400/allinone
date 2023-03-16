package service

import "time"

type ZRTC struct {
	Name string
	Path string
}

type PublishRecord struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	GomssBranch string    `json:"gomss_branch"`
	ZrtcVersion string    `json:"zrtc_version"`
	RecordedAt  time.Time `json:"recorded_at"`
}
