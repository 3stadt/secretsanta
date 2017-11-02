package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/yunabe/easycsv"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
	"time"
)

var auth smtp.Auth

type Config struct {
	SmtpUser     string
	SmtpPass     string
	FromAddress  string
	EmailSubject string
}

type participant struct {
	Name  string `name:"name"`
	Email string `name:"email"`
}

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

type templatedata struct {
	Santa     string
	Presentee string
}

type smtpAuth struct {
	username string
	password string
}

type maildata struct {
	auth         smtpAuth
	subject      string
	templatedata templatedata
}

func main() {
	c, err := readConfig()
	exitOnErr(err)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Path to the CSV file: [participants.csv] ")
	participantsCsv, err := reader.ReadString('\n')
	participantsCsv = strings.TrimSpace(participantsCsv)
	exitOnErr(err)
	if participantsCsv == "" {
		participantsCsv = "participants.csv"
	}
	if _, err := os.Stat(participantsCsv); os.IsNotExist(err) {
		fmt.Printf("file %s does not exist or is not readable\n", participantsCsv)
		os.Exit(1)
	}
	r := easycsv.NewReaderFile(participantsCsv)
	participants := []participant{}
	err = r.ReadAll(&participants)
	exitOnErr(err)
	fmt.Println()
	for _, part := range participants {
		fmt.Printf("%s \t %s\n", part.Name, part.Email)
	}
	fmt.Print("\nE-Mails will be sent to all persons listed above. Continue? [y/N] ")
	do, err := reader.ReadString('\n')
	do = strings.TrimSpace(do)
	exitOnErr(err)
	if strings.ToLower(do) != "y" {
		os.Exit(0)
	}
	pairings := pair(participants)
	m := maildata{
		auth: smtpAuth{
			username: c.SmtpUser,
			password: c.SmtpPass,
		},
		subject: c.EmailSubject,
	}
	for santa, presentee := range pairings {
		m.templatedata = templatedata{
			Santa:     santa.Name,
			Presentee: presentee.Name,
		}
		req := NewRequest([]string{santa.Email}, c.FromAddress, fmt.Sprintf(m.subject, santa.Name), "")
		send(&m, req)
	}
}

func readConfig() (*Config, error) {
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		return nil, errors.New("Please create a config.toml file.")
	}
	tomlBytes, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return nil, err
	}
	tomlData := string(tomlBytes)
	var conf Config
	_, err = toml.Decode(tomlData, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func pair(p []participant) map[participant]participant {
	// shuffle the slice: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle#The_modern_algorithm
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	for i := len(p) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
	// generate pairing
	lastIndex := len(p) - 1
	partMap := make(map[participant]participant)
	for i, part := range p {
		if i == lastIndex {
			partMap[part] = p[0]
			continue
		}
		partMap[part] = p[i+1]
	}
	return partMap
}

func send(m *maildata, r *Request) error {
	// mail sending borrowed from @dhanush: https://gist.github.com/dhanush/f1bac67b659cdd88d3703ea758a313c0
	auth = smtp.PlainAuth("", m.auth.username, m.auth.password, "smtp.gmail.com")
	err := r.ParseTemplate("template.html", m.templatedata)
	if err != nil {
		return err
	}
	_, err = r.SendEmail()
	return err
}

func NewRequest(to []string, from, subject, body string) *Request {
	return &Request{
		from:    from,
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"
	if err := smtp.SendMail(addr, auth, r.from, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
