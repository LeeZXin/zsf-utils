package ginutil

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

var (
	DefaultSuccessResp = BaseResp{
		Code:    0,
		Message: "success",
	}
)

type BaseResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *BaseResp) IsSuccess() bool {
	return r.Code == 0
}

type DataResp[T any] struct {
	BaseResp
	Data T `json:"data"`
}

type PageReq struct {
	Cursor int64 `json:"cursor"`
	Limit  int   `json:"limit"`
}

type PageResp[T any] struct {
	DataResp[[]T]
	Next int64 `json:"next"`
}

type Page2Req struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type Page2Resp[T any] struct {
	DataResp[[]T]
	PageNum    int   `json:"pageNum"`
	TotalCount int64 `json:"totalCount"`
}

func HandleErr(err error, c *gin.Context) {
	if err != nil {
		berr, ok := err.(*bizerr.Err)
		if !ok {
			c.String(http.StatusInternalServerError, "")
		} else {
			c.JSON(http.StatusOK, BaseResp{
				Code:    berr.Code,
				Message: berr.Message,
			})
		}
	}
}

func ShouldBind(obj any, c *gin.Context) bool {
	err := c.ShouldBind(obj)
	if err != nil {
		c.String(http.StatusBadRequest, "")
		return false
	}
	return true
}

func GetClientIp(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

func GetFile(c *gin.Context) (io.ReadCloser, bool, error) {
	contentType := strings.ToLower(c.GetHeader("Content-Type"))
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") ||
		strings.HasPrefix(contentType, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			return nil, false, err
		}
		if c.Request.MultipartForm.File == nil {
			return nil, false, http.ErrMissingFile
		}
		for _, files := range c.Request.MultipartForm.File {
			if len(files) > 0 {
				r, err := files[0].Open()
				return r, true, err
			}
		}
		return nil, false, http.ErrMissingFile
	}
	return c.Request.Body, false, nil
}
