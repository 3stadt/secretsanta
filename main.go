package main

import (
	"bufio"
	"fmt"
	"github.com/3stadt/secretsanta/mail"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/yunabe/easycsv"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	SmtpServer   string
	SmtpPort     int
	SmtpUser     string
	SmtpPass     string
	FromAddress  string
	EmailSubject string
}

type participant struct {
	Name  string `name:"name"`
	Email string `name:"email"`
}

func main() {
	c, err := readConfig()
	exitOnErr(err)
	participants := getParticipants()
	var seed *int64
	if len(os.Args) > 1 {
		seedArg, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err == nil {
			seed = &seedArg
		}
	}
	pairings, seed := pair(participants, seed)
	m := mail.MailData{
		Server:   c.SmtpServer,
		Port:     c.SmtpPort,
		Username: c.SmtpUser,
		Password: c.SmtpPass,
		Subject:  c.EmailSubject,
	}
	for santa, presentee := range pairings {
		m.TemplateData = mail.TemplateData{
			Santa:     santa.Name,
			Presentee: presentee.Name,
			Seed:      seed,
		}
		req := mail.NewRequest([]string{santa.Email}, c.FromAddress, fmt.Sprintf(m.Subject, santa.Name), "")
		err := req.Send(&m)
		if err != nil {
			log.Info(err)
		}
	}
	fmt.Printf("Done. Seed used for pairing was %d, save it and use it as argument to re-send the e-mails from this run. (using the same csv file)\n", *seed)
}

func getParticipants() []participant {
	var participants []participant
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
	return participants
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

func pair(p []participant, seed *int64) (map[participant]participant, *int64) {
	if seed == nil {
		now := time.Now().UnixNano()
		seed = &now
	}
	rand.Seed(*seed)
	perm := rand.Perm(len(p))
	lastIndex := len(perm) - 1
	partMap := make(map[participant]participant)
	for i, randIndex := range perm {
		part := p[randIndex]
		if i == lastIndex {
			partMap[part] = p[perm[0]]
			continue
		}
		partMap[part] = p[perm[i+1]]
	}
	return partMap, seed
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
