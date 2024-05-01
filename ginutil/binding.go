package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func BindParams(c *gin.Context, ptr any) error {
	params := make(map[string]string, len(c.Params))
	for _, param := range c.Params {
		params[param.Key] = param.Value
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  ptr,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(params)
}
