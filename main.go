package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/recoilme/slowpoke"
	"github.com/zserge/webview"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

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
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	defer slowpoke.CloseAll()
	defer ln.Close()

	c := conf{
		host:    "http://" + ln.Addr().String(),
		santaDb: "secretsanta.db",
		confDb:  "config.db",
	}

	go func() {
		http.HandleFunc("/santas", c.handleSantas)
		http.HandleFunc("/css/fonts.css", c.handleFontCss)
		http.HandleFunc("/index.html", c.handleIndexHtml)
		fs := http.FileServer(http.Dir("web"))
		http.Handle("/", http.StripPrefix("/", fs))
		log.Fatal(http.Serve(ln, nil))
	}()

	initialHTML := `<!doctype html>
	<html lang="en">
	<head>
	   <meta charset="UTF-8">
	   <meta name="viewport"
	         content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
	   <meta http-equiv="X-UA-Compatible" content="ie=edge">
		<meta http-equiv="refresh" content="0;url=` + c.host + `/index.html">
	   <title>Document</title>
	</head>
	<body>
	<h1>Starting...</h1>
	</body>
	</html>`

	// TODO create headless mode without GUI
	w := webview.New(webview.Settings{
		URL:       `data:text/html,` + url.PathEscape(initialHTML),
		Width:     800,
		Height:    600,
		Resizable: true,
		Debug:     true,
	})

	w.Run()
}

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

func (c *conf) handleSantas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		allSantas, err := c.getAllSantasAsJson()
		if err != nil {
			log.Printf("%+v", errors.Wrap(err, "could not get all santas"))
			return
		}
		_, err = fmt.Fprint(w, allSantas)
		if err != nil {
			log.Printf("%+v", errors.Wrap(err, "could not write to ResponseWriter"))
			return
		}
	case "POST":
		err := c.saveSanta(r)
		if err != nil {
			log.Printf("%+v", errors.Wrap(err, "could not save santas"))
			return
		}
	}
}

func (c *conf) handleFontCss(w http.ResponseWriter, r *http.Request) {
	t := template.New("fonts.css")
	t, err := t.ParseFiles("./templates/css/fonts.css")
	if err != nil {
		log.Printf("%+v", errors.Wrap(err, "could not parse css font template"))
		return
	}
	w.Header().Set("Content-Type", "text/css")
	err = t.Execute(w, c.host)
	if err != nil {
		log.Printf("%+v", errors.Wrap(err, "could not write to ResponseWriter"))
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
	santas := []santa{}
	keys, err := slowpoke.Keys(c.santaDb, nil, 0, 0, true)
	if err != nil {
		return "", errors.Wrap(err, "could not read Keys from db")
	}
	res := slowpoke.Gets(c.santaDb, keys)
	for k, val := range res {
		if k%2 == 0 {
			continue // this is only the key, not the value
		}
		var s santa
		err = json.Unmarshal(val, &s)
		if err != nil {
			return "", errors.Wrap(err, "could not unmarshal db entry")
		}
		santas = append(santas, s)
	}
	b, err := json.Marshal(santas)
	if err != nil {
		log.Println(err.Error())
		return "", errors.Wrap(err, "could not marshal db entry")
	}
	return string(b), nil
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
