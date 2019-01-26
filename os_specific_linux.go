package main

import (
	"log"
	"os"
	"os/exec"
)

func openExplorer() {
	cmd := exec.Command("xdg-open", "./templates/mail/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
