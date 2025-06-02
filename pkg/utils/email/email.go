package email

import (
	"douyin/config"
	"gopkg.in/mail.v2"
)

// EmailSender 邮件发送器结构体，包含SMTP相关配置信息
type EmailSender struct {
	SmtpHost      string `json:"smtp_host"`       // SMTP服务器地址
	SmtpEmailFrom string `json:"smtp_email_from"` // 发件人邮箱地址
	SmtpPass      string `json:"smtp_pass"`       // 发件人邮箱密码或授权码
}

// NewEmailSender 创建并返回一个新的EmailSender实例，读取配置文件中的邮件相关配置
// NewEmailSender 创建并返回一个新的 EmailSender 实例，读取配置文件中的邮件相关配置
func NewEmailSender() *EmailSender {
	// 确保 GlobalConfig 已经初始化，获取全局配置中的 Email 配置项
	eConfig := config.GlobalConfig.Email // 从全局配置中访问邮件相关的配置信息

	// 检查配置是否正确加载，如果 eConfig 为 nil 则表示配置未加载成功
	if eConfig == nil {
		// 程序直接 panic 并输出错误信息，提醒配置文件中 Email 部分未加载
		panic("Email 配置未加载")
	}

	// 返回一个新的 EmailSender 实例，并将配置中的 SMTP 主机、发送邮箱及密码赋值给对应的字段
	return &EmailSender{
		SmtpHost:      eConfig.SmtpHost,  // SMTP 主机地址：用于发送邮件时的服务器地址
		SmtpEmailFrom: eConfig.SmtpEmail, // SMTP 发送邮箱：作为发件人使用的邮箱地址
		SmtpPass:      eConfig.SmtpPass,  // SMTP 邮箱密码：用于认证的密码或授权码
	}
}

// Send 发送邮件
// 参数说明：
// - data：邮件内容（HTML格式）
// - emailTo：目标接收者的邮箱地址
// - subject：邮件主题
func (s *EmailSender) Send(data, emailTo, subject string) error {
	// 新建邮件消息对象
	m := mail.NewMessage()
	m.SetHeader("From", s.SmtpEmailFrom) // 设置发件人
	m.SetHeader("To", emailTo)           // 设置收件人
	m.SetHeader("Subject", subject)      // 设置邮件主题
	m.SetBody("text/html", data)         // 设置邮件内容（HTML格式）

	// 创建一个邮件拨号器，端口465为SMTP SSL端口
	d := mail.NewDialer(s.SmtpHost, 465, s.SmtpEmailFrom, s.SmtpPass)
	// 强制启用TLS，确保安全传输
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// 发送邮件，如有错误则返回错误信息
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
