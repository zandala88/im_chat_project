package ws

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"im/config"
	"im/model"
	"im/public/etcd"
	"im/public/protocol"
	"im/public/rpc"
	"im/server/cache"
	"im/server/mq"
	"im/service"
	"im/util"
	"time"
)

// GetOutputMsg 组装出下行消息
func GetOutputMsg(cmdType protocol.CmdType, code int32, message proto.Message) ([]byte, error) {
	output := &protocol.Output{
		Type:    cmdType,
		Code:    code,
		CodeMsg: util.GetErrorMessage(int(code)),
		Data:    nil,
	}
	if message != nil {
		msgBytes, err := proto.Marshal(message)
		if err != nil {
			zap.S().Error("[GetOutputMsg] message marshal err:", err)
			return nil, err
		}
		output.Data = msgBytes
	}

	bytes, err := proto.Marshal(output)
	if err != nil {
		zap.S().Error("[GetOutputMsg] output marshal err:", err)
		return nil, err
	}
	return bytes, nil
}

// SendToUser 发送消息到好友
func SendToUser(msg *protocol.Message, userId int64) (int64, error) {
	// 获取接受者 seqId
	seq, err := service.GetUserNextSeq(userId)
	if err != nil {
		zap.S().Error("[消息处理] 获取 seq 失败,err:", err)
		return 0, err
	}
	msg.Seq = seq

	// 发给MQ
	err = mq.MessageMQ.Publish(model.MessageToProtoMarshal(&model.Message{
		UserID:      userId,
		SenderID:    msg.SenderId,
		SessionType: int8(msg.SessionType),
		ReceiverId:  msg.ReceiverId,
		MessageType: int8(msg.MessageType),
		Content:     msg.Content,
		Seq:         seq,
		SendTime:    time.UnixMilli(msg.SendTime),
	}))
	if err != nil {
		zap.S().Error("[消息处理] mq.MessageMQ.Publish(messageBytes) 失败,err:", err)
		return 0, err
	}

	// 如果发给自己的，只落库不进行发送
	if userId == msg.SenderId {
		return seq, nil
	}

	// 组装消息
	bytes, err := GetOutputMsg(protocol.CmdType_CT_Message, util.Ok, &protocol.PushMsg{Msg: msg})
	if err != nil {
		zap.S().Error("[消息处理] GetOutputMsg Marshal error,err:", err)
		return 0, err
	}

	// 进行推送
	return 0, Send(userId, bytes)
}

// Send 消息转发
// 是否在线 ---否---> 不进行推送
//    |
//    是
//    ↓
//  是否在本地 --否--> RPC 调用
//    |
//    是
//    ↓
//  消息发送

func Send(receiverId int64, bytes []byte) error {
	// 查询是否在线
	rpcAddr, err := cache.GetUserOnline(receiverId)
	if err != nil {
		return err
	}

	// 不在线
	if rpcAddr == "" {
		zap.S().Debug("[消息处理]，用户不在线，receiverId:", receiverId)
		return nil
	}

	zap.S().Info("[消息处理] 用户在线，rpcAddr:", rpcAddr)

	// 查询是否在本地
	conn := ConnManager.GetConn(receiverId)
	if conn != nil {
		// 发送本地消息
		conn.SendMsg(receiverId, bytes)
		zap.S().Debug("[消息处理]， 发送本地消息给用户, ", receiverId)
		return nil
	}

	// rpc 调用
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = rpc.GetServerClient(rpcAddr).DeliverMessage(ctx, &protocol.DeliverMessageReq{
		ReceiverId: receiverId,
		Data:       bytes,
	})

	if err != nil {
		zap.S().Error("[消息处理] DeliverMessage err, err:", err)
		return err
	}

	return nil
}

// SendToGroup 发送消息到群
func SendToGroup(msg *protocol.Message) error {
	// 获取群成员信息
	userIds, err := service.GetGroupUser(msg.ReceiverId)
	if err != nil {
		zap.S().Error("[群聊消息处理] 查询失败，err:", err, msg)
		return err
	}

	// userId set
	m := make(map[int64]struct{}, len(userIds))
	for _, userId := range userIds {
		m[userId] = struct{}{}
	}

	// 检查当前用户是否属于该群
	if _, ok := m[msg.SenderId]; !ok {
		zap.S().Debug("[群聊消息处理] 用户不属于该群组，msg:", msg)
		return nil
	}

	// 自己不再进行推送
	delete(m, msg.SenderId)

	sendUserIds := make([]int64, 0, len(m))
	for userId := range m {
		sendUserIds = append(sendUserIds, userId)
	}
	// 批量获取 seqId
	seqs, err := service.GetUserNextSeqBatch(sendUserIds)
	if err != nil {
		zap.S().Error("[批量获取 seq] 失败,err:", err)
		return err
	}

	//  k:userid v:该userId的seq
	sendUserSet := make(map[int64]int64, len(seqs))
	for i, userId := range sendUserIds {
		sendUserSet[userId] = seqs[i]
	}

	// 创建 Message 对象
	messages := make([]*model.Message, 0, len(m))
	for userId, seq := range sendUserSet {
		messages = append(messages, &model.Message{
			UserID:      userId,
			SenderID:    msg.SenderId,
			SessionType: int8(msg.SessionType),
			ReceiverId:  msg.ReceiverId,
			MessageType: int8(msg.MessageType),
			Content:     msg.Content,
			Seq:         seq,
			SendTime:    time.UnixMilli(msg.SendTime),
		})
	}

	// 发给MQ
	err = mq.MessageMQ.Publish(model.MessageToProtoMarshal(messages...))
	if err != nil {
		zap.S().Error("[消息处理] 群聊消息发送 MQ 失败,err:", err)
		return err
	}

	// 组装消息，进行推送
	userId2Msg := make(map[int64][]byte, len(m))
	for userId, seq := range sendUserSet {
		msg.Seq = seq
		bytes, err := GetOutputMsg(protocol.CmdType_CT_Message, util.Ok, &protocol.PushMsg{Msg: msg})
		if err != nil {
			zap.S().Error("[消息处理] GetOutputMsg Marshal error,err:", err)
			return err
		}
		userId2Msg[userId] = bytes
	}

	// 获取全部网关服务，进行消息推送
	services := etcd.DiscoverySer.GetServices()
	local := fmt.Sprintf("%s:%s", config.Configs.App.IP, config.Configs.App.RPCPort)
	for _, addr := range services {
		// 如果是本机，进行本地推送
		if local == addr {
			zap.S().Debug("进行本地推送")
			GetServer().SendMessageAll(userId2Msg)
		} else {
			zap.S().Debug("远端推送：", addr)
			// 如果不是本机，进行远程 RPC 调用
			_, err = rpc.GetServerClient(addr).DeliverMessageAll(context.Background(), &protocol.DeliverMessageAllReq{
				ReceiverId_2Data: userId2Msg,
			})

			if err != nil {
				zap.S().Error("[消息处理] DeliverMessageAll err, err:", err)
				return err
			}
		}
	}

	return nil
}
