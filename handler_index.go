package main

import (
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

func (c *conf) handleIndexHtml(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html")
	t, err := t.ParseFiles("./templates/index.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, c.Host)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
