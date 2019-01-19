package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/3stadt/secretsanta/mail"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/recoilme/slowpoke"
	log "github.com/sirupsen/logrus"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type mailConf struct {
	Server   string
	Port     string
	Username string
	Password string
	Subject  string
}

type conf struct {
	host    string
	santaDb string
	confDb  string
}

type santa struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
}

func main() {
	//file, err := os.OpenFile("secretsanta.log", os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.SetOutput(file)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	f, err := os.Create("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	e := toml.NewEncoder(w)
	err = e.Encode(mail.MailData{
		Server:   "Server",
		Port:     9090,
		Username: "User",
		Password: "Pass",
		Subject:  "Subj",
	})
	if err != nil {
		log.Fatal(err)
	}

	//ln, err := net.Listen("tcp", "127.0.0.1:0")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer slowpoke.CloseAll()
	//defer ln.Close()
	//
	//c := conf{
	//	host:    "http://" + ln.Addr().String(),
	//	santaDb: "secretsanta.db",
	//	confDb:  "config.db",
	//}
	//
	//go func() {
	//	r := mux.NewRouter()
	//	r.HandleFunc("/santas", c.handlePostSanta).Methods("POST")
	//	r.HandleFunc("/santas", c.handleGetSanta).Methods("GET")
	//	r.HandleFunc("/santas/{mail}", c.handleDeleteSanta).Methods("DELETE")
	//	r.HandleFunc("/css/fonts.css", c.handleFontCss).Methods("GET")
	//	r.HandleFunc("/index.html", c.handleIndexHtml).Methods("GET")
	//	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web"))))
	//	log.Error(http.Serve(ln, r))
	//}()
	//
	//initialHTML := `<!doctype html>
	//<html lang="en">
	//<head>
	//   <meta charset="UTF-8">
	//   <meta name="viewport"
	//         content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
	//   <meta http-equiv="X-UA-Compatible" content="ie=edge">
	//	<meta http-equiv="refresh" content="0;url=` + c.host + `/index.html">
	//   <title>Document</title>
	//</head>
	//<body>
	//<h1>Starting...</h1>
	//</body>
	//</html>`
	//
	//// TODO create headless mode without GUI
	//w := webview.New(webview.Settings{
	//	URL:       `data:text/html,` + url.PathEscape(initialHTML),
	//	Width:     800,
	//	Height:    600,
	//	Resizable: true,
	//	Debug:     true,
	//})
	//
	//w.Run()
}

//func (c *conf) handleSendMail(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	var seed *int64 = nil
//	if _, ok := vars["seed"]; ok {
//		userSeed, err := strconv.ParseInt(vars["seed"], 10, 64)
//		if err != nil {
//			w.WriteHeader(http.StatusBadRequest)
//			_, _ = fmt.Fprint(w, err.Error())
//			return
//		}
//		seed = &userSeed
//	}
//	santas, err := c.getAllSantas()
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//	pairings, seed := pair(santas, seed)
//	m := mail.MailData{
//		Server:   c.SmtpServer,
//		Port:     c.SmtpPort,
//		Username: c.SmtpUser,
//		Password: c.SmtpPass,
//		Subject:  c.EmailSubject,
//	}
//	for santa, presentee := range pairings {
//		m.TemplateData = mail.TemplateData{
//			Santa:     santa.Name,
//			Presentee: presentee.Name,
//			Seed:      seed,
//		}
//		req := mail.NewRequest([]string{santa.Email}, c.FromAddress, fmt.Sprintf(m.Subject, santa.Name), "")
//		err := req.Send(&m)
//		if err != nil {
//			log.Info(err)
//		}
//	}
//}

func (c *conf) handleIndexHtml(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html")
	t, err := t.ParseFiles("./templates/index.html")
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

func (c *conf) handleGetSanta(w http.ResponseWriter, r *http.Request) {
	allSantas, err := c.getAllSantasAsJson()
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not get all santas"))
		return
	}
	_, err = fmt.Fprint(w, allSantas)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not write to ResponseWriter"))
		return
	}
}

func (c *conf) handlePostSanta(w http.ResponseWriter, r *http.Request) {
	err := c.saveSanta(r)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not save santas"))
		return
	}
}

func (c *conf) handleDeleteSanta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := slowpoke.Delete(c.santaDb, []byte(vars["mail"]))
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "could not delete santa"))
		return
	}
	return
}

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

func (c *conf) saveSanta(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return errors.Wrap(err, "could not parse form")
	}
	newSanta := santa{
		Name: r.FormValue("name"),
		Mail: r.FormValue("mail"),
	}
	b, _ := json.Marshal(newSanta)
	err := slowpoke.Set(c.santaDb, []byte(newSanta.Mail), b)
	if err != nil {
		return errors.Wrap(err, "could not write to db")
	}
	return nil
}

func (c *conf) getAllSantasAsJson() (string, error) {
	santas, err := c.getAllSantas()
	if err != nil {
		return "", errors.Wrap(err, "could not marshal db entry")
	}
	b, err := json.Marshal(santas)
	if err != nil {
		return "", errors.Wrap(err, "could not marshal db entry")
	}
	return string(b), nil
}

func (c *conf) getAllSantas() ([]santa, error) {
	santas := []santa{}
	keys, err := slowpoke.Keys(c.santaDb, nil, 0, 0, true)
	if err != nil {
		return santas, errors.Wrap(err, "could not read Keys from db")
	}
	res := slowpoke.Gets(c.santaDb, keys)
	for k, val := range res {
		if k%2 == 0 {
			continue // this is only the key, not the value
		}
		var s santa
		err = json.Unmarshal(val, &s)
		if err != nil {
			return santas, errors.Wrap(err, "could not unmarshal db entry")
		}
		santas = append(santas, s)
	}
	return santas, nil
}

func pair(p []santa, seed *int64) (map[santa]santa, *int64) {
	if seed == nil {
		now := time.Now().UnixNano()
		seed = &now
	}
	rand.Seed(*seed)
	perm := rand.Perm(len(p))
	lastIndex := len(perm) - 1
	partMap := make(map[santa]santa)
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

func generateRandString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
