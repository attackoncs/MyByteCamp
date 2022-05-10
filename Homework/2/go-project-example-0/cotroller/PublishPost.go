package cotroller

import (
	"github.com/Moonlight-Zhao/go-project-example/service"
	"strconv"
)

func PublishPost(parent_idStr, content string) *PageData {

	parent_id, err := strconv.ParseInt(parent_idStr, 10, 64)
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}

	pageInfo, err := service.PublishPost(parent_id, content)
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	return &PageData{
		Code: 0,
		Msg:  "success",
		Data: pageInfo,
	}

}
