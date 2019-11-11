package main

import "github.com/3stadt/secretsanta/mail"

type conf struct {
	Host     string
	santaDb  string
	confFile string
	MailData *mail.Data
}

type santa struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
}

const initialHTML = `<!doctype html>
	<html lang="en">
	<head>
	  <meta charset="UTF-8">
	  <meta name="viewport"
	        content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
	  <meta http-equiv="X-UA-Compatible" content="ie=edge">
		<meta http-equiv="refresh" content="0;url={{ . }}/index.html">
	  <title>Document</title>
	</head>
	<body>
	<h1>Starting...</h1>
	</body>
	</html>`
