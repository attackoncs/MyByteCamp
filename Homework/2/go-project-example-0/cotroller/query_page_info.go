package cotroller

import (
	"strconv"

	"github.com/Moonlight-Zhao/go-project-example/service"
)

// PageData 最终发送给客户端的json数据对应的结构体，我们需要错误码，以及对应错误码对应的消息，最后再是数据(用空接口实现泛型
type PageData struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// QueryPageINfo 真正和客户端进行交互的函数，需要注意客户端发来的流量都是字符串形式
//查找主题topicIdStr的回帖列表
func QueryPageInfo(topicIdStr string) *PageData {
	topicId, err := strconv.ParseInt(topicIdStr, 10, 64)
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	pageInfo, err := service.QueryPageInfo(topicId)
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
