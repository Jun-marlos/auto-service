package mail

import (
	"context"
	"fmt"
	"testing"
)

func TestMail(t *testing.T) {
	// 邮件接收方
	mailTo := []string{
		//可以是多个接收人
		"919897089@qq.com",
	}

	subject := "Hello World!" // 邮件主题
	body := "测试发送邮件"          // 邮件正文

	err := SendMail(mailTo, subject, body)
	if err != nil {
		fmt.Println("Send fail! - ", err)
		return
	}
	fmt.Println("Send successfully!")
}

func TestVerifyCode(t *testing.T) {
	ctx := context.Background()
	token, _ := SendVerifyCode(ctx, "919897089@qq.com")
	fmt.Println(token)
}
