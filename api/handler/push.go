package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"im_chat_project/api/rpc"
	"im_chat_project/proto"
	"im_chat_project/public"
	"im_chat_project/public/config"
	"strconv"
)

type FormPush struct {
	Msg       string `form:"msg" json:"msg" binding:"required"`
	ToUserId  string `form:"toUserId" json:"toUserId" binding:"required"`
	RoomId    int64  `form:"roomId" json:"roomId" binding:"required"`
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func Push(c *gin.Context) {
	var formPush FormPush
	if err := c.ShouldBindBodyWith(&formPush, binding.JSON); err != nil {
		public.ResponseError(c, 40002, err.Error())
		return
	}
	authToken := formPush.AuthToken
	msg := formPush.Msg
	toUserId := formPush.ToUserId
	toUserIdInt, _ := strconv.ParseInt(toUserId, 10, 64)
	getUserNameReq := &proto.GetUserInfoRequest{UserId: toUserIdInt}
	code, toUserName := rpc.RpcLogicObj.GetUserNameByUserId(getUserNameReq)
	if code == public.CodeFail {
		public.ResponseError(c, 40003, "rpc fail get friend userName")
		return
	}
	checkAuthReq := &proto.CheckAuthRequest{AuthToken: authToken}
	code, fromUserId, fromUserName := rpc.RpcLogicObj.CheckAuth(checkAuthReq)
	if code == public.CodeFail {
		public.ResponseError(c, 40004, "rpc fail get self info")
		return
	}
	roomId := formPush.RoomId
	req := &proto.Send{
		Msg:          msg,
		FromUserId:   fromUserId,
		FromUserName: fromUserName,
		ToUserId:     toUserIdInt,
		ToUserName:   toUserName,
		RoomId:       roomId,
		Op:           config.OpSingleSend,
	}
	code, rpcMsg := rpc.RpcLogicObj.Push(req)
	if code == public.CodeFail {
		public.ResponseError(c, 40005, rpcMsg)
		return
	}
	public.ResponseSuccess(c, nil)
	return
}

type FormRoom struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
	Msg       string `form:"msg" json:"msg" binding:"required"`
	RoomId    int64  `form:"roomId" json:"roomId" binding:"required"`
}

func PushRoom(c *gin.Context) {
	var formRoom FormRoom
	if err := c.ShouldBindBodyWith(&formRoom, binding.JSON); err != nil {
		public.ResponseError(c, 40002, err.Error())
		return
	}
	authToken := formRoom.AuthToken
	msg := formRoom.Msg
	roomId := formRoom.RoomId
	checkAuthReq := &proto.CheckAuthRequest{AuthToken: authToken}
	authCode, fromUserId, fromUserName := rpc.RpcLogicObj.CheckAuth(checkAuthReq)
	if authCode == public.CodeFail {
		public.ResponseError(c, 40003, "rpc fail get self info")
		return
	}
	req := &proto.Send{
		Msg:          msg,
		FromUserId:   fromUserId,
		FromUserName: fromUserName,
		RoomId:       roomId,
		Op:           config.OpRoomSend,
	}
	code, msg := rpc.RpcLogicObj.PushRoom(req)
	if code == public.CodeFail {
		public.ResponseError(c, 40006, "rpc push room msg fail!")
		return
	}
	public.ResponseSuccess(c, msg)
	return
}

type FormCount struct {
	RoomId int64 `form:"roomId" json:"roomId" binding:"required"`
}

func Count(c *gin.Context) {
	var formCount FormCount
	if err := c.ShouldBindBodyWith(&formCount, binding.JSON); err != nil {
		public.ResponseError(c, 40002, err.Error())
		return
	}
	roomId := formCount.RoomId
	req := &proto.Send{
		RoomId: roomId,
		Op:     config.OpRoomCountSend,
	}
	code, msg := rpc.RpcLogicObj.Count(req)
	if code == public.CodeFail {
		public.ResponseError(c, 40007, "rpc get room count fail!")
		return
	}
	public.ResponseSuccess(c, msg)
	return
}

type FormRoomInfo struct {
	RoomId int64 `form:"roomId" json:"roomId" binding:"required"`
}

func GetRoomInfo(c *gin.Context) {
	var formRoomInfo FormRoomInfo
	if err := c.ShouldBindBodyWith(&formRoomInfo, binding.JSON); err != nil {
		public.ResponseError(c, 40002, err.Error())
		return
	}
	roomId := formRoomInfo.RoomId
	req := &proto.Send{
		RoomId: roomId,
		Op:     config.OpRoomInfoSend,
	}
	code, msg := rpc.RpcLogicObj.GetRoomInfo(req)
	if code == public.CodeFail {
		public.ResponseError(c, 40008, "rpc get room info fail!")
		return
	}
	public.ResponseSuccess(c, msg)
	return
}
