package protocol

import (
	"encoding/base64"
	"google.golang.org/protobuf/proto"
	"testing"
)

func Test_BuildProtocol(t *testing.T) {
	loginMsg := &LoginMsg{
		Token: []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MzA4NzIyMzcsImlzcyI6InphbmRhbGEifQ.JbfqTMz4qWvAf0G-rsE46Eu-xxAitRv3oC0C31a9-vs"),
	}
	tmp, _ := proto.Marshal(loginMsg)

	input := &Input{
		Type: CmdType_CT_Login,
		Data: tmp,
	}
	marshal, err := proto.Marshal(input)
	if err != nil {
		t.Error(err)
		return
	}

	// 将二进制数据编码为 Base64
	base64Str := base64.StdEncoding.EncodeToString(marshal)
	t.Log(base64Str)
}
