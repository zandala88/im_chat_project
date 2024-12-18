package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"im/public/protocol"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

var phoneNumbers = []string{
	"16797178986", "16536807149", "19198567094", "14378244943", "17779153238", "19407368684", "18461589678", "19937796807", "15685158537",
	"19415526542", "14858131428", "19962946509", "14207648261", "13884172772", "17301459721", "18765410507", "15937792460", "13898488660",
	"18996907734", "15910740413", "19555564583", "15640017686", "14619888526", "18573158497", "17188800934", "14644624176", "17400854016",
	"13396856280", "17020892378", "17958474806", "15985200924", "19750365135", "13688625988", "16185813144", "14221697332", "19335513038",
	"18677383730", "17348666660", "16361146536", "17595886009", "13471114162", "15598385358", "16817670718", "19441976651", "13000694534",
	"18667157140", "18914106918", "17095516305", "16409030801", "16558826053", "13551133627", "13742599771", "15560066743", "16954101110",
	"13336333952", "14736832229", "14694635862", "16309586963", "13779246067", "16568870243", "14872483297", "14710219980", "19339809611",
	"16874650178", "19251432779", "14986066545", "13470606065", "13940614382", "17240644867", "13148427196", "16646991136", "15872536348",
	"13949815654", "13940785096", "18582298506", "16730946643", "15933656103", "15589207464", "17596817656", "15862161645", "17688964439",
	"13825557657", "18878089362", "15457532973", "18221515017", "19350091483", "15905416066", "16218637836", "16991347430", "15530989907",
	"19271195768", "13820825011", "19911697456", "17054069543", "14794761535", "19700854148", "17635332250", "18952287880", "18989259020",
	"14175628859", "14789968118", "18226198384", "15367998154", "19277086631", "16028544370", "16558085625", "18435014053", "15305750064",
	"16774790467", "15847962936", "18621142211", "17447043846", "18385522260", "16966075439", "19650586408", "16325244458", "17307558206",
	"15963727259", "15434459789", "19366242263", "14674639557", "13556184817", "18557869086", "19266887662", "18101177170", "18534397183",
	"14522871033", "14246857047", "13290680594", "15605725524", "18730448320", "14201428926", "18880242429", "18913512007", "18842955024",
	"19030990239", "14555781441", "18508755434", "17883368028", "18470831854", "19979181449", "16342802781", "16731654450", "17908908023",
	"19570861006", "15785470465", "15620539844", "13993170104", "17603730057", "18940276190", "15480575625", "15809020787", "19256727440",
	"16240350586", "16763547277", "19261664902", "14407694648", "14172825689", "13053367278", "17743335382", "13400278102", "13999865070",
	"18536330379", "13881144918", "14098353767", "15036648185", "13280342714", "13345857859", "15411691432", "14766423168", "13373926836",
	"16484500439", "18382153056", "18314458816", "13707167501", "14196041803", "15663247240", "13404088434", "19485415182", "18979831574",
	"19599424479", "17727301332", "17422537711", "13177085821", "15456568224", "13192296690", "19616951705", "17206719783", "14080964265",
	"16256293890", "16155870994", "15275099041", "15242350016", "18901950892", "16275943013", "19930516666", "17224062663", "14351603716",
	"16987164991", "15155532126",
}

const (
	loginURL       = "http://localhost:9090/login"
	websocketURL   = "ws://localhost:9091/ws" // 替换为实际 WebSocket 地址
	password       = "123456"
	concurrentUser = 10 // 并发用户数
)

type LoginResponse struct {
	Code int `json:"code"`
	Data struct {
		Token  string `json:"token"`
		UserID string `json:"user_id"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func login(phoneNumber string, password string) (string, error) {
	// 构造表单数据
	formData := fmt.Sprintf("phone_number=%s&password=%s", phoneNumber, password)
	req, err := http.NewRequest("POST", loginURL, bytes.NewBufferString(formData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed: status %d", resp.StatusCode)
	}

	// 解析登录响应
	var loginResp LoginResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return "", err
	}

	if loginResp.Code != 200 {
		return "", fmt.Errorf("login failed: %s", loginResp.Msg)
	}

	return loginResp.Data.Token, nil
}

func connectWebSocket(token string) error {
	// 使用 token 设置 WebSocket 请求头
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// 建立 WebSocket 连接
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:9091/ws", header) // WebSocket 地址根据实际修改
	if err != nil {
		return err
	}

	// 准备发送的消息
	loginMsg := &protocol.LoginMsg{
		Token: []byte(token),
	}

	// 将 LoginMsg 序列化为二进制格式
	tmp, err := proto.Marshal(loginMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal login message: %v", err)
	}

	// 包装 Input 消息
	input := &protocol.Input{
		Type: protocol.CmdType_CT_Login, // 假设 `CmdType_CT_Login` 是你定义的命令类型
		Data: tmp,
	}

	// 序列化 Input 消息
	marshal, err := proto.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input message: %v", err)
	}

	// 发送二进制消息到 WebSocket
	err = conn.WriteMessage(websocket.BinaryMessage, marshal)
	if err != nil {
		return fmt.Errorf("failed to send message over WebSocket: %v", err)
	}

	log.Println("Message sent successfully over WebSocket")
	go send(conn)
	return nil
}

func send(conn *websocket.Conn) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		msg := &protocol.UpMsg{
			ClientId: 1,
			Msg: &protocol.Message{
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

		input := &protocol.Input{
			Type: protocol.CmdType_CT_Message,
			Data: tmp,
		}
		marshal, err := proto.Marshal(input)
		if err != nil {
			fmt.Errorf("failed to marshal input message: %v", err)
			return
		}
		// 发送二进制消息到 WebSocket
		err = conn.WriteMessage(websocket.BinaryMessage, marshal)
		if err != nil {
			fmt.Errorf("failed to send message over WebSocket: %v", err)
		}
	}
}

func Test_perf(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(len(phoneNumbers))

	for i := 0; i < len(phoneNumbers); i++ {
		go func(phone string) {
			defer wg.Done()
			token, err := login(phone, password)
			if err != nil {
				log.Printf("Login failed for phone %s: %v\n", phone, err)
				return
			}
			log.Printf("Login successful for phone %s, token: %s\n", phone, token)

			go func(t string) {
				err := connectWebSocket(t)
				if err != nil {
					log.Printf("WebSocket connection failed for token %s: %v", t, err)
				}
			}(token)
		}(phoneNumbers[i])
	}

	wg.Wait()

	select {}
}
