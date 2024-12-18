package main

import (
	"flag"
)

/*
------ 目标：单机群聊压测 ------
注：本机压本机
系统：Windows 10 19045.2604
CPU: 2.90 GHz AMD Ryzen 7 4800H with Radeon Graphics
内存: 16.0 GB
群成员总数: 500人
在线人数: 500人
每秒/次发送消息数量: 500条
每秒理论响应消息数量：25 0000条 = 500条 * 在线500人
发送消息次数: 40次
响应消息总量：1000 0000条 = 500条 * 在线500人 * 40次
Message 表数量总数：1000 0000条 = 总数500人 * 500条 * 40次
丢失消息数量: 0条
总耗时: 39948ms
平均每500条消息发送/转发在线人员/在线人员接收总耗时: 998ms（其实更短，因为消息是每秒发一次）

如果发送消息次数为 1，时间为：940ms
*/

var (
	// 起始phone_num
	pn = flag.Int64("pn", 100000, "First phone num")
	// 群成员总数
	gn = flag.Int64("gn", 200, "群成员总数")
	// 在线成员数量
	on = flag.Int64("on", 200, "在线成员数量")
	// 每次发送消息数量
	sn = flag.Int64("sn", 200, "每次发送消息数量")
	// 发送消息次数
	tn = flag.Int64("tn", 1, "发送消息次数")
)

func main() {
	flag.Parse()

	mgr := NewManager(*pn, *on, *sn, *gn, *tn)
	mgr.Run()

	select {}
}
