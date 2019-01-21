package main

import (
	"fmt"
	"github.com/3stadt/secretsanta/mail"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (c *conf) handleGetMailConfig(w http.ResponseWriter, r *http.Request) {

}

func (c *conf) handlePostMailConfig(w http.ResponseWriter, r *http.Request) {

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
		c.mailData.TemplateData = mail.TemplateData{
			Santa:     santa.Name,
			Presentee: presentee.Name,
			Seed:      seed,
		}
		req := mail.NewRequest([]string{santa.Mail}, c.mailData.FromAddress, fmt.Sprintf(c.mailData.Subject, santa.Name), "")
		err := req.Send(c.mailData)
		if err != nil {
			log.Errorf("%+v", errors.Wrap(err, "could not send mail"))
			return
		}
	}
}
