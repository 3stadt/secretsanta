package main

import (
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

func (c *conf) handlePreviewHtml(w http.ResponseWriter, r *http.Request) {
	t := template.New("preview.html")
	t, err := t.ParseFiles("./templates/preview.html")
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, c.host)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
