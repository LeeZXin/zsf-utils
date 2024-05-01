package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	err = decoder.Decode(params)
	if err != nil {
		return err
	}
	return validate(ptr)
}

func BindQuery(c *gin.Context, ptr any) error {
	err := binding.MapFormWithTag(ptr, c.Request.URL.Query(), "json")
	if err != nil {
		return err
	}
	return validate(ptr)
}

func validate(obj any) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
