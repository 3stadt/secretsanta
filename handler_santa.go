package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/recoilme/slowpoke"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (c *conf) handleGetSanta(w http.ResponseWriter, r *http.Request) {
	allSantas, err := c.getAllSantasAsJson()
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not get all santas"))
		return
	}
	_, err = fmt.Fprint(w, allSantas)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not write to ResponseWriter"))
		return
	}
}

func (c *conf) handlePostSanta(w http.ResponseWriter, r *http.Request) {
	err := c.saveSanta(r)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not save santas"))
		return
	}
}

func (c *conf) handleDeleteSanta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := slowpoke.Delete(c.santaDb, []byte(vars["mail"]))
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not delete santa"))
		return
	}
	return
}
