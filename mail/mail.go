package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
	"net/textproto"
)

type Mail interface {
	Send(m *Data) error
}

type Pairing struct {
	Santa     string
	Presentee string
	Seed      *int64
}

type TemplateData struct {
	Headline      string `json:"headline"`
	GreetingIntro string `json:"greetingIntro"`
	Santa         string `json:"santa"`
	GreetingOutro string `json:"greetingOutro"`
	SantaMatch    string `json:"santaMatch"`
	Intro         string `json:"intro"`
	Outro         string `json:"outro"`
	Greeting      string `json:"greeting"`
}

type Data struct {
	Server       string `json:"smtpServer"`
	Port         int    `json:"smtpPort"`
	Username     string `json:"smtpUser"`
	Password     string `json:"smtpPass"`
	Subject      string `json:"mailSubject"`
	FromAddress  string `json:"senderAddress"`
	SSL          bool   `json:"smtpSsl"`
	Pairing      Pairing
	TemplateData *TemplateData
}

type mailReq struct {
	MailData *Data
	Request  *request
}

type request struct {
	from    string
	to      []string
	subject string
	body    string
	ssl     bool
}

func NewRequest(to []string, from, subject, body string, ssl bool) *request {
	return &request{
		from:    from,
		to:      to,
		subject: subject,
		body:    body,
		ssl:     ssl,
	}
}

func (r *request) Send(m *Data, templatePath string) error {
	err := r.parseTemplate(templatePath, m.TemplateData)
	if err != nil {
		return err
	}
	d := mailReq{m, r}
	return d.sendViaSmtp()
}

func (r *mailReq) sendViaSmtp() error {
	e := &email.Email{
		To:      r.Request.to,
		From:    r.Request.from,
		Subject: r.Request.subject,
		HTML:    []byte(r.Request.body),
		Headers: textproto.MIMEHeader{},
	}
	if r.Request.ssl {
		return e.SendWithTLS(
			fmt.Sprintf("%s:%d", r.MailData.Server, r.MailData.Port),
			smtp.PlainAuth("", r.MailData.Username, r.MailData.Password, r.MailData.Server),
			&tls.Config{
				ServerName: r.MailData.Server,
			},
		)
	}
	return e.Send(fmt.Sprintf("%s:%d", r.MailData.Server, r.MailData.Port), smtp.PlainAuth("", r.MailData.Username, r.MailData.Password, r.MailData.Server))
}

func (r *request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
