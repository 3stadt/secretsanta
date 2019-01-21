package mail

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
)

type Mail interface {
	Send(m *MailData) error
}

type TemplateData struct {
	Santa     string
	Presentee string
	Seed      *int64
}

type MailData struct {
	Server       string
	Port         int
	Username     string
	Password     string
	Subject      string
	FromAddress string
	TemplateData TemplateData
}

type mailReq struct {
	MailData *MailData
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

func (r *request) Send(m *MailData) error {
	err := r.parseTemplate("template.html", m.TemplateData)
	if err != nil {
		return err
	}
	d := mailReq{m, r}
	return d.sendViaSmtp()
}

func (r *mailReq) sendViaSmtp() error {
	m := gomail.NewMessage()
	m.SetHeader("From", r.Request.from)
	m.SetHeader("To", r.Request.to...)
	m.SetHeader("Subject", r.Request.subject)
	m.SetBody("text/html", r.Request.body)
	d := gomail.NewPlainDialer(r.MailData.Server, r.MailData.Port, r.MailData.Username, r.MailData.Password)
	return d.DialAndSend(m)
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
