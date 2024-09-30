package main

import (
	"encoding/base64"
	"github.com/gogo/protobuf/proto"
	"im/internal/websocket/protocol"
	"testing"
)

func Test_Protobuf(t *testing.T) {
	msg := &protocol.Message{
		To:          "1c19ddda-e510-440c-b495-989a58a87e1f",
		From:        "4b0c3427-fdb1-4874-8df6-029151dd595a",
		MessageType: 1,
		ContentType: 1,
		Content:     "test",
	}
	ans, err := proto.Marshal(msg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ans)
	base64String := base64.StdEncoding.EncodeToString(ans)
	t.Log(base64String)
	//tmpMsg := &protocol.Message{}
	//err = proto.Unmarshal(ans, tmpMsg)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Logf("%#v", tmpMsg)
}
