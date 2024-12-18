### 消息类型

`protobuf` 序列化后发送



### 连接方式

`ws://localhost:9091/ws`



### 消息格式

##### 统一的最外层

客户端发送前：先组装出下层消息例如 HeartBeatMsg，序列化作为 Input 的 data 值，再填写 type 值，序列化 Input 发送给服务端

##### 客户端到服务端的消息最外层

```
// type 类型 
// 1 连接注册，客户端向服务端发送，建立连接
// 2 心跳，客户端向服务端发送，连接保活
// 3 消息投递，可能是服务端发给客户端，也可能是客户端发给服务端
// 4 ACK
// 5 离线消息同步
{
	type: 1, // 消息类型
	data: "" // 内嵌结构序列化
}
```

##### 服务端到客户端的消息最外层

```
{
	type: 1,// 类型同上共用
	code :200,// 错误码 200 正常
	codeMsg: "",// 错误码信息
	data:"", // 内嵌结构序列化 
}
```





##### type = 1

连接上websocket后，发送该类型消息，用作登录

等待服务端回复ACK信息，超时则重传

````
客户端到服务端
{
	token："" // 接口登录时返回的token
}

服务端到客户端的ACK
回复外层type = 4 AND 内层 type = 3 
````



##### type = 2

心跳检测机制

暂不设定内容，type = 2 对于即可

```
{

}
```



##### type = 3

客户端 发送该消息类型时，每次都将clientId++，并启动一个计时器，检测服务端是否返回（ACK类型消息，type = 4）(**仍未实现**)，若无则再次发送该消息

服务端处理后同样恢复ACK，具体格式见下方

```
// 发送消息
{
	clientId: 1 // 
	msg : {
		session_type: 1, // 1 单聊类型消息 2 群聊类型消息
		receiver_id： 1, // 接收者id 用户id/群组id
		sender_id: 1, // 发送者id
		message_type: 1,// 1 文字类型消息 2 图片类型消息 3 语音类型消息
		content："", // 实际内容
		seq: 1,// 客户端的最大消息同步序号 用作离线消息同步
		send_time： 312312 // 消息发送时间戳，ms
	}
}

// 服务端的ACK消息
回复外层type = 4 AND 内层 type = 1 


// 所发送的对象接收消息
{ 
	msg : {
		session_type: 1, // 1 单聊类型消息 2 群聊类型消息
		receiver_id： 1, // 接收者id 用户id/群组id
		sender_id: 1, // 发送者id
		message_type: 1,// 1 文字类型消息 2 图片类型消息 3 语音类型消息
		content："", // 实际内容
		seq: 1,// 客户端的最大消息同步序号 用作离线消息同步
		send_time： 312312 // 消息发送时间戳，ms
	}
}

// 客户端回复ACK
回复外层type = 4 AND 内层 type = 2 
```



##### type = 4

ACK消息类型，用作确认消息到达

客户端主动发送场景：收到由服务端发送的非ACK类型消息后，回复ACK

服务端发送场景：

- 收到登录消息回复对应类型消息ACK
- 收到待转送的消息 回复对应类型消息ACK

```
{
	type: 1,// 1 服务端回复客户端发来的消息，2 客户端回复服务端发来的消息 3 登录
	clientId:1,//
	seq: 1,// 
}
```



##### type = 5

离线消息同步

```
{
	seq: 3 // 
}

{
	messages : [
		{
			session_type: 1, // 1 单聊类型消息 2 群聊类型消息
			receiver_id： 1, // 接收者id 用户id/群组id
			sender_id: 1, // 发送者id
			message_type: 1,// 1 文字类型消息 2 图片类型消息 3 语音类型消息
			content："", // 实际内容
			seq: 1,// 客户端的最大消息同步序号 用作离线消息同步
			send_time： 312312 // 消息发送时间戳，ms
		},
	],
	has_more: false // 是否还有更多
}
```



