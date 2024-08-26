package public

import (
	"context"
	"fmt"
	"github.com/jordan-wright/email"
	"go.uber.org/zap"
	"im/config"
	"im/util"
	"math/rand"
	"net/smtp"
	"regexp"
	"strconv"
	"time"
)

func SendConfirmEmail(targetMail string) int {
	// 检查是否在3分钟内发送过邮件
	if Redis.Get(context.Background(), "send-email:"+targetMail).Val() != "" {
		zap.S().Error("send email too frequently")
		return util.SendEmailIn3Min
	}

	code := getConfirmCode()
	Redis.Set(context.Background(), "email:"+targetMail, code, time.Minute*5)

	err := sendConfirmMessage(targetMail, code)
	if err != nil {
		zap.S().Errorf("SendConfirmEmail : %v", err)
		return util.SendEmailError
	}

	Redis.Set(context.Background(), "send-email:"+targetMail, code, time.Minute*3)

	return 0
}

func VerifyEmailFormat(email string) bool {
	pattern := `^[^\s@]+@[^\s@]+\.[^\s@]+$` //match email
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func sendConfirmMessage(targetMail string, code string) error {
	em := email.NewEmail()
	em.From = fmt.Sprintf("zandala-im <%s>", config.Configs.Email.Addr)
	em.To = []string{targetMail}

	// 邮箱标题
	em.Subject = "Email Confirm Code "

	// 邮件内容
	emailContent := "你的验证码为 " + code + ", Your code will expire in 5 minutes"
	em.Text = []byte(emailContent)

	// 调用发送邮件函数
	err := em.Send(config.Configs.Email.Service+":25", smtp.PlainAuth("", config.Configs.Email.Addr, config.Configs.Email.Key, config.Configs.Email.Service))
	if err != nil {
		zap.S().Errorf("SendConfirmMessage : %v", err)
		return err
	}

	return nil
}

func getConfirmCode() string {
	var confirmCode int
	for i := 0; i < 6; i++ {
		confirmCode = confirmCode*10 + (rand.Intn(9) + 1) //随机函数获取值
	}
	// 转换成字符串
	confirmCodeStr := strconv.Itoa(confirmCode)
	return confirmCodeStr
}
