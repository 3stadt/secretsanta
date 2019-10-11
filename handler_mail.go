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
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, c)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (c *conf) handleSendMail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var seed *int64 = nil
	if _, ok := vars["seed"]; ok {
		userSeed, err := strconv.ParseInt(vars["seed"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, err.Error())
			return
		}
		seed = &userSeed
	}
	santas, err := c.getAllSantas()
	if err != nil {
		log.Println(err.Error())
		return
	}
	pairings, seed := pair(santas, seed)

	for santa, presentee := range pairings {
		c.MailData.Pairings = mail.Pairings{
			Santa:     santa.Name,
			Presentee: presentee.Name,
			Seed:      seed,
		}
		req := mail.NewRequest([]string{santa.Mail}, c.MailData.FromAddress, fmt.Sprintf(c.MailData.Subject, santa.Name), "")
		err := req.Send(c.MailData)
		if err != nil {
			log.Errorf("%+v", errors.Wrap(err, "could not send mail"))
			return
		}
	}
}
