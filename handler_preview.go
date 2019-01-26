package main

import (
	"encoding/json"
	"fmt"
	"github.com/3stadt/secretsanta/explorer"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"path/filepath"
)

func (c *conf) handleOpenTemplateFolder(w http.ResponseWriter, r *http.Request) {
	explorer.Open()
}

func (c *conf) handlePreviewList(w http.ResponseWriter, r *http.Request) {
	files, err := filepath.Glob("./templates/mail/*.html")
	if err != nil {
		log.Fatal(err)
	}
	for i, file := range files {
		files[i] = filepath.Base(file)
	}
	jsonFileList, err := json.Marshal(files)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = fmt.Fprintf(w, string(jsonFileList))
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (c *conf) handlePreviewMail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := filepath.Base(vars["filename"]) // prevent changing folder via relative path

	t := template.New(filename)
	t, err := t.ParseFiles(fmt.Sprintf("./templates/mail/%s", filename))
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := struct {
		Santa      string
		SantaMatch string
	}{
		"Danielle",
		"Benjamin",
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (c *conf) handlePreviewPage(w http.ResponseWriter, r *http.Request) {
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
