package main

import (
	"encoding/json"
	"fmt"
	"github.com/recoilme/slowpoke"
	"github.com/zserge/webview"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
)

type conf struct {
	host string
	db   string
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
		host: "http://" + ln.Addr().String(),
		db:   "secretsanta.db",
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
		santas := []santa{}
		keys, err := slowpoke.Keys(c.db, nil, 0, 0, true)
		if err != nil {
			log.Println(err.Error())
			return
		}
		res := slowpoke.Gets(c.db, keys)
		for k, val := range res {
			if k%2 == 0 {
				continue // this is only the key, not the value
			}
			var s santa
			err = json.Unmarshal(val, &s)
			if err != nil {
				log.Println(err.Error())
				return
			}
			santas = append(santas, s)
		}
		b, err := json.Marshal(santas)
		if err != nil {
			log.Println(err.Error())
			return
		}
		_, err = fmt.Fprint(w, string(b))
		if err != nil {
			log.Println(err.Error())
			return
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm() err: %v", err)
			return
		}
		newSanta := santa{
			Name: r.FormValue("name"),
			Mail: r.FormValue("mail"),
		}
		b, _ := json.Marshal(newSanta)
		err := slowpoke.Set(c.db, []byte(newSanta.Mail), b)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (c *conf) handleFontCss(w http.ResponseWriter, r *http.Request) {
	t := template.New("fonts.css")
	t, err := t.ParseFiles("./templates/css/fonts.css")
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/css")
	err = t.Execute(w, c.host)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
