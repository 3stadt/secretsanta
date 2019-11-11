package explorer

import (
	"log"
	"os"
	"os/exec"
)

func Open() {
	cmd := exec.Command("start", "./templates/mail/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
