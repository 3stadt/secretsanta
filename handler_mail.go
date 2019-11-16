package main

import (
	"fmt"
	"github.com/3stadt/secretsanta/mail"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strconv"
)

func (c *conf) handleConfigPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("config.html")
	t, err := t.ParseFiles("./templates/config.html")
	if err != nil {
		log.Error(err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, c)
	if err != nil {
		log.Error(err)
		return
	}
}

func (c *conf) handleConfigPagePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error(err)
		return
	}
	c.MailData.Server = r.FormValue("smtpServer")
	port, err := strconv.ParseInt(r.FormValue("smtpPort"), 10, 0)
	if err != nil {
		log.Error(err)
		port = 25
	}
	c.MailData.Port = int(port)
	c.MailData.Username = r.FormValue("smtpUser")
	c.MailData.Password = r.FormValue("smtpPass")
	c.MailData.Subject = r.FormValue("mailSubject")
	c.MailData.FromAddress = r.FormValue("senderAddress")
	ssl, err := strconv.ParseBool(r.FormValue("smtpSsl"))
	if err != nil {
		log.Error(err)
		ssl = false
	}
	c.MailData.SSL = ssl
	err = c.MailData.WriteConfig(c.confFile)
	if err != nil {
		log.Error(err)
		return
	}
	c.sendMail(w, r)
}

func (c *conf) sendMail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var seed *int64 = nil
	if _, ok := vars["seed"]; ok {
		userSeed, err := strconv.ParseInt(vars["seed"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, err.Error())
			log.Error(err)
			return
		}
		seed = &userSeed
	}
	santas, err := c.getAllSantas()
	if err != nil {
		log.Error(err)
		return
	}
	pairings, seed := pair(santas, seed)

	for santa, presentee := range pairings {
		c.MailData.Pairing = mail.Pairing{
			Santa:     santa.Name,
			Presentee: presentee.Name,
			Seed:      seed,
		}
		c.MailData.TemplateData.Santa = santa.Name
		c.MailData.TemplateData.SantaMatch = presentee.Name
		req := mail.NewRequest([]string{santa.Mail}, c.MailData.FromAddress, c.MailData.Subject, "", c.MailData.SSL)
		err := req.Send(c.MailData, "templates/mail/001.html")
		if err != nil {
			log.Errorf("%+v", errors.Wrap(err, "could not send mail"))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, fmt.Sprintf("%+v", errors.Wrap(err, "could not send mail")))
			return
		}
	}
}
