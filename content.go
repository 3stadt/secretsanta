package main

import (
	"encoding/json"
	"github.com/3stadt/secretsanta/mail"
	"github.com/pkg/errors"
	"github.com/recoilme/slowpoke"
	"net/http"
)

func (c *conf) saveContent(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return errors.Wrap(err, "could not parse form")
	}
	content := mail.TemplateData{
		Headline:      r.FormValue("headline"),
		GreetingIntro: r.FormValue("greetingIntro"),
		GreetingOutro: r.FormValue("greetingOutro"),
		Intro:         r.FormValue("intro"),
		Outro:         r.FormValue("outro"),
		Greeting:      r.FormValue("greeting"),
	}
	b, _ := json.Marshal(content)
	err := slowpoke.Set(c.santaDb, []byte("mailContent"), b)
	if err != nil {
		return errors.Wrap(err, "could not write to db")
	}
	return nil
}

func (c *conf) getContent() error {
	data := mail.TemplateData{}
	bytes, err := slowpoke.Get(c.santaDb, []byte("mailContent"))
	if err != nil {
		return errors.Wrap(err, "could not read content from db")
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}
	if c.MailData == nil {
		c.MailData = &mail.Data{
			Server:       "",
			Port:         0,
			Username:     "",
			Password:     "",
			Subject:      "",
			FromAddress:  "",
			Pairings:     mail.Pairings{},
			TemplateData: &data,
		}
		return nil
	}
	c.MailData.TemplateData = &data
	return nil
}
