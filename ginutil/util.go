package ginutil

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/gin-gonic/gin"
	"net/http"
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
