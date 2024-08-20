package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// 定义一个返回 json 的结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 失败的响应
func (resp *Response) Fail(response http.ResponseWriter) {
	resp.Code = -1
	res, err := json.Marshal(resp)
	if err != nil {
		log.Println(err.Error())
	}
	//response.Header().Set("Access-Control-Allow-Origin", "*")
	//response.Header().Set("Access-Control-Allow-Credentials", "true")
	//response.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	//response.Header().Set("Access-Control-Allow-Headers", "DNT,Keep-Alive,User-Agent,Cache-Control,Content-Type,Authorization,x-token")
	//// 设置 header 为 json, 返回 json
	//response.Header().Set("Content-Encoding","zh-CN")
	//response.Header().Set("Content-Type", "application/json")
	// 设置响应状态 200
	response.WriteHeader(http.StatusOK)
	response.Write(res)
}

// 成功的响应
func (resp *Response) Success(response http.ResponseWriter) {
	resp.Code = 0
	res, err := json.Marshal(resp)
	if err != nil {
		log.Println(err.Error())
	}
	//response.Header().Set("Access-Control-Allow-Origin", "*")
	//response.Header().Set("Access-Control-Allow-Credentials", "true")
	//response.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	//response.Header().Set("Access-Control-Allow-Headers", "DNT,Keep-Alive,User-Agent,Cache-Control,Content-Type,Authorization,x-token")
	//// 设置 header 为 json, 返回 json
	//response.Header().Set("Content-Encoding","zh-CN")
	//response.Header().Set("Content-Type", "application/json")
	// 设置响应状态 200
	response.WriteHeader(http.StatusOK)
	response.Write(res)
}
