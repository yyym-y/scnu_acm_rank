package middle

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/smtp"
	"scnu_acm_rank/biz/config"

	"github.com/jordan-wright/email"
)

type emailConfig struct {
	from     string
	subject  string
	password string
	host     string
}

var E emailConfig

func init() {
	E = emailConfig{}
	config.Add(&E)
}

func (e *emailConfig) Update() {
	e.host = config.Conf.EmailHost
	e.password = config.Conf.EmailPassword
	e.from = config.Conf.EmailFrom
	e.subject = config.Conf.EmailSubject
}

func SendEmail(to []string) error {
	e := email.NewEmail()
	e.From = E.from
	e.To = to
	code := ""
	for i := 0; i < 6; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10))
	}
	e.Subject = E.subject
	e.HTML = []byte("欢迎注册scnu_acm_rank, 您的验证码是: <h1>" + code + "</h1> ")
	err := e.SendWithStartTLS("smtp.163.com:25", smtp.PlainAuth("", E.from, E.password, "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		return err
	}
	AddCode(to[0], code)
	return nil
}

func SendKeyEmail(to []string, code string) error {
	e := email.NewEmail()
	e.From = E.from
	e.To = to
	e.Subject = E.subject
	e.HTML = []byte("您所创建队伍的口令为： <h1>" + code + "</h1>")
	err := e.SendWithStartTLS("smtp.163.com:25", smtp.PlainAuth("", E.from, E.password, "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		return err
	}
	AddCode(to[0], code)
	return nil
}
