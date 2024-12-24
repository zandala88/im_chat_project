package protocol

import (
	"encoding/base64"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func Test_BuildProtocol(t *testing.T) {
	loginMsg := &LoginMsg{
		Token: []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJleHAiOjE3MzcwMTM2OTEsImlzcyI6InphbmRhbGEifQ.KlHtkxoDSI9tsNW1kb96GvpIweXwaLaQdNuh26sw8JE"),
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

func Test_SaveMsg(t *testing.T) {
	msg := &UpMsg{
		ClientId: 1,
		Msg: &Message{
			SessionType: 1,
			ReceiverId:  2,
			SenderId:    3,
			MessageType: 1,
			Content:     []byte("hello"),
			Seq:         1,
			SendTime:    time.Now().Unix(),
		},
	}

	tmp, _ := proto.Marshal(msg)

	input := &Input{
		Type: CmdType_CT_Message,
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
	base64Str := "CAQQyAEaB1N1Y2Nlc3MiBggBEAEYAQ=="
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

func Test_PushProtocol2(t *testing.T) {
	base64Str := "CAM="
	decoded, _ := base64.StdEncoding.DecodeString(base64Str)
	output := &ACKMsg{}
	err := proto.Unmarshal(decoded, output)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(output)
	t.Log(output.ClientId)
	t.Log(output.Seq)
}

func Test_PushProtocol3(t *testing.T) {
	base64Str := "CAMSIQodCAEQAhjIAyABKgQxMTExMNmzj7++MjjZs4+/vjIQAQ=="
	decoded, _ := base64.StdEncoding.DecodeString(base64Str)
	output := &Input{}
	err := proto.Unmarshal(decoded, output)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(output)
	t.Log(output.Type)
	t.Logf("%#v", output.Data)
	output2 := &UpMsg{}
	err = proto.Unmarshal(output.Data, output2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(output2)
}

func Test_PushProtocol4(t *testing.T) {
	base64Str := "Cg0IARACGAMgASoBMTAC"
	decoded, _ := base64.StdEncoding.DecodeString(base64Str)
	output := &SyncOutputMsg{}
	err := proto.Unmarshal(decoded, output)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(output)
	t.Log(output.Messages)
	t.Log(output.HasMore)
}
