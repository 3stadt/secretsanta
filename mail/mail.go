package mail

import (
	"net/smtp"
	"bytes"
	"html/template"
)

type Mail interface {
	Send(m *MailData) error
}

type TemplateData struct {
	Santa     string
	Presentee string
}

type SmtpAuth struct {
	Username string
	Password string
}

type MailData struct {
	Auth         SmtpAuth
	Subject      string
	TemplateData TemplateData
}

type request struct {
	from    string
	to      []string
	subject string
	body    string
	auth    smtp.Auth
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
	r.auth = smtp.PlainAuth("", m.Auth.Username, m.Auth.Password, "smtp.gmail.com")
	err := r.parseTemplate("template.html", m.TemplateData)
	if err != nil {
		return err
	}
	_, err = r.sendViaSmtp()
	return err
}

func (r *request) sendViaSmtp() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"
	if err := smtp.SendMail(addr, r.auth, r.from, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
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
