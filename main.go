package main

import (
	"bytes"
	"github.com/3stadt/secretsanta/mail"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/recoilme/slowpoke"
	log "github.com/sirupsen/logrus"
	"github.com/zserge/webview"
	"html/template"
	"net"
	"net/http"
	"net/url"
)

func main() {
	//file, err := os.OpenFile("secretsanta.log", os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.SetOutput(file)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	c := conf{
		santaDb:  "secretsanta.db",
		confFile: "config.toml",
	}

	mc, err := mail.ReadConfig(c.confFile)
	if err != nil {
		log.Info(errors.Wrap(err, "could not read mail config"))
	}
	c.mailData = mc

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	defer slowpoke.CloseAll()
	defer ln.Close()

	c.host = "http://" + ln.Addr().String()

	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/mailconfig", c.handlePostMailConfig).Methods("POST")
		r.HandleFunc("/mailconfig", c.handleGetMailConfig).Methods("GET")
		r.HandleFunc("/santas", c.handlePostSanta).Methods("POST")
		r.HandleFunc("/santas", c.handleGetSanta).Methods("GET")
		r.HandleFunc("/santas/{mail}", c.handleDeleteSanta).Methods("DELETE")
		r.HandleFunc("/css/fonts.css", c.handleFontCss).Methods("GET")
		r.HandleFunc("/index.html", c.handleIndexHtml).Methods("GET")
		r.HandleFunc("/preview.html", c.handlePreviewHtml).Methods("GET")
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web"))))
		log.Error(http.Serve(ln, r))
	}()

	html := new(bytes.Buffer)
	t := template.Must(template.New("initialHtml").Parse(initialHTML))
	err = t.Execute(html, c.host)
	if err != nil {
		log.Fatal(err)
	}

	w := webview.New(webview.Settings{
		URL:       `data:text/html,` + url.PathEscape(html.String()),
		Width:     800,
		Height:    600,
		Resizable: true,
		Debug:     true,
	})

	w.Run()
}
