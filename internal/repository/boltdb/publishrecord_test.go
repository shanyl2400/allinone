package boltdb

import (
	"gomssbuilder/internal/repository/model"
	"testing"
)

var (
	records = []*model.PublishRecord{
		{
			Version:     "v1",
			GomssBranch: "v1",
			ZrtcVersion: "zrtcv1",
		},
		{
			Version:     "v2",
			GomssBranch: "v2",
			ZrtcVersion: "zrtcv2",
		},
		{
			Version:     "v3",
			GomssBranch: "v3",
			ZrtcVersion: "zrtcv3",
		},
		{
			Version:     "v4",
			GomssBranch: "v4",
			ZrtcVersion: "zrtcv4",
		},
	}
)

func TestPutView(t *testing.T) {
	err := GetClient().Open()
	if err != nil {
		t.Fatal(err)
	}

	repo, err := NewBoltdbPublishRecordRepository()
	if err != nil {
		t.Fatal(err)
	}

	for _, record := range records {
		err = repo.Put(record)
		if err != nil {
			t.Error(err)
		}
	}

	output, err := repo.Recent(5)
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range output {
		t.Logf("%#v", r)
	}
}
