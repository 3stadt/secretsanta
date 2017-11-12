package main

import (
	"github.com/hoisie/web"
	"github.com/monoculum/formam"
	"github.com/phayes/freeport"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/matcornic/hermes"
	"github.com/3stadt/secretsanta/mail"
)

var (
	beatCount = 0
	logger    = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
)

type Participant struct {
	Name  string
	Email string
}

type Formdata struct {
	SmtpUser     string
	SmtpPass     string
	SmtpFrom     string
	SmtpServer   string
	SmtpPort     int
	Subject      string
	Participants []Participant
	MailContent  string
	Seed         string
}

func main() {
	port := getFreePort()
	s := web.NewServer()

	s.SetLogger(logger)
	s.Config = &web.ServerConfig{
		StaticDir: "/home/n/go/src/github.com/3stadt/secretsanta/assets/docroot", // TODO make dynamic
	}
	s.Post("/api/heartbeat", heartbeat)
	s.Post("/api/sendmail", formendpoint)
	showBrowser("http://127.0.0.1:" + port)
	if len(os.Args) < 2 || os.Args[1] != "dev" {
		go checkHeartbeat(s)
	}
	s.Run("127.0.0.1:" + port)
}

func showBrowser(url string) {
	err := open.Run(url)
	if err != nil {
		logger.Printf("unable to open browser: %s\n", err.Error())
		os.Exit(1)
	}
}

func checkHeartbeat(s *web.Server) {
	time.Sleep(5 * time.Second) // 10 seconds grace time for first opening the app
	for {
		time.Sleep(5 * time.Second)
		if beatCount > 0 {
			logger.Printf("found heartbeat. [%d]\n", beatCount)
			beatCount = 0
			continue
		}
		s.Close()
		logger.Print("lost heartbeat from browser, exiting...\n")
		os.Exit(0)
	}
}

func heartbeat() string {
	beatCount++
	return ""
}

func formendpoint(ctx *web.Context) string {
	r := ctx.Request
	r.ParseForm()
	fd := Formdata{}
	dec := formam.NewDecoder(&formam.DecoderOptions{})
	if err := dec.Decode(r.Form, &fd); err != nil {
		return err.Error()
	}
	html, err := buildMail(fd.MailContent)
	if err != nil {
		return err.Error()
	}
	fd.MailContent = html
	return sendMail(&fd).Error()
}

func sendMail(fd *Formdata) error {
	var participants []string
	md := &mail.MailData{
		Auth: mail.SmtpAuth{
			Username: fd.SmtpUser,
			Password: fd.SmtpPass,
		},
	}
	for _, p := range fd.Participants {
		participants = append(participants, p.Name+"<"+p.Email+">")
	}
	m := mail.NewRequest(participants, fd.SmtpFrom, fd.Subject, fd.MailContent)
	return m.Send(md)
}

func buildMail(c string) (string, error) {
	h := hermes.Hermes{
		Theme: new(SecretSantaTheme),
		Product: hermes.Product{
			Name: "Secret Santa",
		},
	}
	email := hermes.Email{
		Body: hermes.Body{
			FreeMarkdown: hermes.Markdown(c),
		},
	}
	return h.GenerateHTML(email)
}

func getFreePort() string {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	return strconv.Itoa(port)
}
