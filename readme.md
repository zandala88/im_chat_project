### swagger 

http://127.0.0.1:8080/swagger/index.html


### 消息格式

```json
{
	"userId": 0, // 自己的ID
	"cmd": 0,  // 1 私聊 2 群聊 
	"toId": 0,	// 发送对象ID
	"media": 0,  // 消息样式（预留字段）
	"content": "demoString", // 消息内容
	"pic": "demoString", // 图片（预留）
	"url": "demoString", // 预留 
	"memo": "demoString", // 消息描述
	"amount": 0, // 预留
}

```