package main

import (
	"github.com/zserge/webview"
	"log"
	"net"
	"net/http"
	"net/url"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	host := "http://" + ln.Addr().String()
	go func() {
		http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {

		})
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
		<meta http-equiv="refresh" content="0;url=` + host + `">
	   <title>Document</title>
	</head>
	<body>
	<h1>Starting...</h1>
	</body>
	</html>`

	//err = webview.Open("SecretSanta", "", 800, 600, false)

	w := webview.New(webview.Settings{
		URL:       `data:text/html,` + url.PathEscape(initialHTML),
		Width:     800,
		Height:    600,
		Resizable: true,
		Debug:     true,
	})

	w.Run()
}
