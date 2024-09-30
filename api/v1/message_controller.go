package v1

import (
	"go.uber.org/zap"
	"im/internal/model"
	"im/internal/service"
	"im/util"

	"github.com/gin-gonic/gin"
)

// 获取消息列表
func GetMessage(c *gin.Context) {
	messageRequest := &model.MessageRequest{}
	err := c.BindQuery(&messageRequest)
	if nil != err {
		zap.S().Error("bindQuery err = ", err)
	}
	zap.S().Infof("messageRequest params: %#v", messageRequest)

	messages, err := service.MessageService.GetMessages(messageRequest)
	if err != nil {
		util.FailResp(c, err.Error())
		return
	}

	util.SuccessResp(c, messages)
}
