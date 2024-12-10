package protocol

import (
	"encoding/base64"
	"google.golang.org/protobuf/proto"
	"testing"
)

func Test_BuildProtocol(t *testing.T) {
	loginMsg := &LoginMsg{
		Token: []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJleHAiOjE3MzYzMzA1NTIsImlzcyI6InphbmRhbGEifQ.Sr-69_uaJ6IptxYE0PDaiDqXW8uXW6DvzHNjBhUF40I"),
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

func Test_PushProtocol(t *testing.T) {
	base64Str := "ChIIBBDIARoHU3VjY2VzcyICCAM="
	decoded, _ := base64.StdEncoding.DecodeString(base64Str)
	t.Logf("%q", decoded)
	output := &Output{}
	err := proto.Unmarshal(decoded, output)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(output)
	t.Logf("Unmarshaled Output: %#v", output)
	t.Log(output.Data)
	t.Log(output.Code)
}
