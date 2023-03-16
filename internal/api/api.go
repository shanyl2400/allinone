package api

import (
	"gomssbuilder/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type API struct {
	builder *service.GomssBuilder
}

func (a *API) getGomssBranches(c *gin.Context) {
	branches, err := a.builder.ListGomssBranches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errResponse{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, listBranchResponse{
		Branches: branches,
	})
}

func (a *API) getZRTCPath(c *gin.Context) {
	zrtcs, err := a.builder.ListZrtc()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errResponse{
			Message: err.Error(),
		})
		return
	}
	ans := make([]*zrtc, 0, len(zrtcs))
	for _, item := range zrtcs {
		ans = append(ans, &zrtc{
			Name: item.Name,
			Path: item.Path,
		})
	}

	c.JSON(http.StatusOK, listZRTCResopnse{
		Zrtcs: ans,
	})
}

func (a *API) getPublishLogs(c *gin.Context) {
	logs := a.builder.PublishLogs()
	c.JSON(http.StatusOK, getPublishLogsResonse{
		Logs: logs,
	})
}

func (a *API) getRecentPublish(c *gin.Context) {
	records, err := a.builder.RecentPublish()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errResponse{
			Message: err.Error(),
		})
		return
	}
	ans := make([]*publishRecord, 0, len(records))
	for _, item := range records {
		ans = append(ans, &publishRecord{
			ID:          item.ID,
			Version:     item.Version,
			GomssBranch: item.GomssBranch,
			ZrtcVersion: item.ZrtcVersion,
			RecordedAt:  item.RecordedAt.UnixMilli(),
		})
	}

	c.JSON(http.StatusOK, listPublishRecordersResponse{
		Records: ans,
	})
}

func (a *API) publish(c *gin.Context) {
	req := new(publishRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &errResponse{
			Message: err.Error(),
		})
		return
	}
	err = a.builder.Publish(req.GomssBranch, req.ZRTCPath, req.Version, req.LocalZRTC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errResponse{
			Message: err.Error(),
		})
		return
	}
	c.JSON(200, &errResponse{
		Message: "success",
	})
}

func newAPI() *API {
	return &API{
		builder: service.NewBuilder(),
	}
}
