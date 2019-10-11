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
	"os"
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
	c.MailData = mc

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	defer slowpoke.CloseAll()
	defer ln.Close()

	c.Host = "http://" + ln.Addr().String()

	go func() {
		r := mux.NewRouter()
		// Create a new Santa
		r.HandleFunc("/santas", c.handlePostSanta).Methods("POST")
		// Get all Santas as JSON
		r.HandleFunc("/santas", c.handleGetSanta).Methods("GET")
		// Delete a Santa from DB, using the mail address as identifier
		r.HandleFunc("/santas/{mail}", c.handleDeleteSanta).Methods("DELETE")
		// Send the actual mail. Returns an error if no config is saved yet
		r.HandleFunc("/mail/send", c.handleSendMail).Methods("POST")
		// The mail preview is embedded in the preview.html file, this endpoint shows the actual mail and is not the preview page
		r.HandleFunc("/mail/template/{filename}", c.handlePreviewMail).Methods("GET")
		// Return all filenames from HTML files inside the template/mail folder
		r.HandleFunc("/previews/available", c.handlePreviewList).Methods("GET")
		// Save the mail content/text to the DB
		r.HandleFunc("/content", c.handlePostMailContent).Methods("POST")
		// Convenience function: Open the OS Explorer at the template folder
		r.HandleFunc("/os/openExplorer", c.handleOpenTemplateFolder).Methods("GET")

		r.HandleFunc("/css/fonts.css", c.handleFontCss).Methods("GET")
		r.HandleFunc("/index.html", c.handleIndexHtml).Methods("GET")
		r.HandleFunc("/content.html", c.handleMailContent).Methods("GET")
		r.HandleFunc("/preview.html", c.handlePreviewPage).Methods("GET")
		r.HandleFunc("/config.html", c.handleConfigPage).Methods("GET")
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web"))))
		log.Error(http.Serve(ln, r))
	}()

	html := new(bytes.Buffer)
	t := template.Must(template.New("initialHtml").Parse(initialHTML))
	err = t.Execute(html, c.Host)
	if err != nil {
		log.Fatal(err)
	}

	w := webview.New(webview.Settings{
		Title:     "The secret santa matcher",
		URL:       `data:text/html,` + url.PathEscape(html.String()),
		Width:     800,
		Height:    900,
		Resizable: true,
		Debug:     os.Getenv("DEBUG") == "true",
	})

	w.Run()
}
