package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"im_chat_project/api/rpc"
	"im_chat_project/proto"
	"im_chat_project/public"
)

type FormLogin struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"passWord" json:"passWord" binding:"required"`
}

func Login(c *gin.Context) {
	var formLogin FormLogin
	if err := c.ShouldBindBodyWith(&formLogin, binding.JSON); err != nil {
		public.ResponseError(c, 40002, err.Error())
		return
	}
	req := &proto.LoginRequest{
		Name:     formLogin.UserName,
		Password: public.Sha1(formLogin.Password),
	}
	code, authToken, msg := rpc.RpcLogicObj.Login(req)
	if code == public.CodeFail || authToken == "" {
		public.ResponseError(c, 40008, msg)
		return
	}
	public.ResponseSuccess(c, authToken)
}

type FormRegister struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"passWord" json:"passWord" binding:"required"`
}

func Register(c *gin.Context) {
	var formRegister FormRegister
	if err := c.ShouldBindBodyWith(&formRegister, binding.JSON); err != nil {
		public.ResponseError(c, 40002, err.Error())
		return
	}
	req := &proto.RegisterRequest{
		Name:     formRegister.UserName,
		Password: public.Sha1(formRegister.Password),
	}
	code, authToken, msg := rpc.RpcLogicObj.Register(req)
	if code == public.CodeFail || authToken == "" {
		public.ResponseError(c, 40009, msg)
		return
	}
	public.ResponseSuccess(c, authToken)
}

type FormCheckAuth struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func CheckAuth(c *gin.Context) {
	var formCheckAuth FormCheckAuth
	if err := c.ShouldBindBodyWith(&formCheckAuth, binding.JSON); err != nil {
		public.ResponseError(c, 40002, err.Error())
		return
	}
	authToken := formCheckAuth.AuthToken
	req := &proto.CheckAuthRequest{
		AuthToken: authToken,
	}
	code, userId, userName := rpc.RpcLogicObj.CheckAuth(req)
	if code == public.CodeFail {
		public.ResponseError(c, 40010, "auth fail")
		return
	}
	var jsonData = map[string]interface{}{
		"userId":   userId,
		"userName": userName,
	}
	public.ResponseSuccess(c, jsonData)
}

type FormLogout struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func Logout(c *gin.Context) {
	var formLogout FormLogout
	if err := c.ShouldBindBodyWith(&formLogout, binding.JSON); err != nil {
		public.ResponseError(c, 40002, err.Error())
		return
	}
	authToken := formLogout.AuthToken
	logoutReq := &proto.LogoutRequest{
		AuthToken: authToken,
	}
	code := rpc.RpcLogicObj.Logout(logoutReq)
	if code == public.CodeFail {
		public.ResponseError(c, 40011, "logout fail!")
		return
	}
	public.ResponseSuccess(c, nil)
}
