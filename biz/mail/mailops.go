package mail

import (
	"context"
	"time"

	"github.com/cloudwego/goapi/biz/dal"
	util "github.com/cloudwego/goapi/biz/utils"
	"gopkg.in/gomail.v2"
)

/*
go邮件发送
*/

func SendMail(mailTo []string, subject string, body string) error {
	m := gomail.NewMessage(
		//发送文本时设置编码，防止乱码。 如果txt文本设置了之后还是乱码，那可以将原txt文本在保存时
		//就选择utf-8格式保存
		gomail.SetEncoding(gomail.Base64),
	)
	m.SetHeader("From", m.FormatAddress(MAIL_ADDRESS, "Auto-Life-Z")) // 添加别名
	m.SetHeader("To", mailTo...)                                      // 发送给用户(可以多个)
	m.SetHeader("Subject", subject)                                   // 设置邮件主题
	m.SetBody("text/html", body)                                      // 设置邮件正文

	/*
	   创建SMTP客户端，连接到远程的邮件服务器，需要指定服务器地址、端口号、用户名、密码，如果端口号为465的话，
	   自动开启SSL，这个时候需要指定TLSConfig
	*/
	d := gomail.NewDialer(MAIL_HOST, MAIL_PORT, MAIL_ADDRESS, MAIL_PWD) // 设置邮件正文
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(m)
	return err
}

func SendVerifyCode(ctx context.Context, email string) (string, error) {
	token := util.RandomStringCreate()
	code := util.RandomCodeCreate()
	keytoken2code := "[verify_code]" + token
	keytoken2mail := "[token2email]" + token
	err := dal.RedisAdd(ctx, keytoken2code, code, time.Minute*5)
	if err != nil {
		return "", err
	}
	err = dal.RedisAdd(ctx, keytoken2mail, email, time.Hour*2)
	if err != nil {
		return "", err
	}
	mailTo := []string{
		//可以是多个接收人
		email,
	}
	err = SendMail(mailTo, VERIFY_SUBJECT, VERIFY_BODY+code)
	if err != nil {
		return "", err
	}
	return token, nil
}
