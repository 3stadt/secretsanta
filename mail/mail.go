package mail

import (
	"bytes"
	"fmt"
	"github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
	"net/textproto"
)

type Mail interface {
	Send(m *Data) error
}

type Pairings struct {
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
	Server       string
	Port         int
	Username     string
	Password     string
	Subject      string
	FromAddress  string
	Pairings     Pairings
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
}

func NewRequest(to []string, from, subject, body string) *request {
	return &request{
		from:    from,
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *request) Send(m *Data) error {
	err := r.parseTemplate("template.html", m.Pairings)
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
