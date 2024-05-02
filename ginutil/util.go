package ginutil

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PageDirection string

const (
	PreviousPage PageDirection = "previous"
	NextPage     PageDirection = "next"
)

func (d PageDirection) IsValid() bool {
	switch d {
	case PreviousPage, NextPage:
		return true
	default:
		return false
	}
}

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
	Cursor    int64         `json:"cursor"`
	Limit     int           `json:"limit"`
	Direction PageDirection `json:"direction"`
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
