package main

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

func (c *conf) handleMailContent(w http.ResponseWriter, r *http.Request) {
	t := template.New("content.html")
	t, err := t.ParseFiles("./templates/content.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = c.getContent()
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = t.Execute(w, c)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (c *conf) handlePostMailContent(w http.ResponseWriter, r *http.Request) {
	err := c.saveContent(r)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not save content"))
		return
	}
}
