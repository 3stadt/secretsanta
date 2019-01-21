package main

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

func (c *conf) handleFontCss(w http.ResponseWriter, r *http.Request) {
	t := template.New("fonts.css")
	t, err := t.ParseFiles("./templates/css/fonts.css")
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not parse css font template"))
		return
	}
	w.Header().Set("Content-Type", "text/css")
	err = t.Execute(w, c.host)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not write to ResponseWriter"))
		return
	}
}
