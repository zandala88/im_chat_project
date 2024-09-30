package v1

import (
	"go.uber.org/zap"
	"im/config"
	"im/internal/service"
	"im/util"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 前端通过文件名称获取文件流，显示文件
func GetFile(c *gin.Context) {
	fileName := c.Param("fileName")
	zap.S().Info("GetFile", fileName)
	data, _ := os.ReadFile(config.Configs.StaticPath.FilePath + fileName)
	c.Writer.Write(data)
}

// 上传头像等文件
func SaveFile(c *gin.Context) {
	namePreffix := uuid.New().String()

	userUuid := c.PostForm("uuid")

	file, _ := c.FormFile("file")
	fileName := file.Filename
	index := strings.LastIndex(fileName, ".")
	suffix := fileName[index:]

	newFileName := namePreffix + suffix

	zap.S().Info("file", zap.String("file name", config.Configs.StaticPath.FilePath+newFileName))
	zap.S().Info("userUuid", zap.String("userUuid name", userUuid))

	c.SaveUploadedFile(file, config.Configs.StaticPath.FilePath+newFileName)
	err := service.UserService.ModifyUserAvatar(newFileName, userUuid)
	if err != nil {
		util.FailResp(c, err.Error())
	}

	util.SuccessResp(c, newFileName)
}
