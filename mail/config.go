package mail

import (
	"bufio"
	"encoding/base64"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func ReadConfig(configFile string) (*MailData, error) {
	var mailConf MailData
	log.Infof("reading config file %s", configFile)
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(b), &mailConf); err != nil {
		return nil, err
	}
	pwByte, err := base64.StdEncoding.DecodeString(mailConf.Password)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode base64 password")
	}
	mailConf.Password = string(pwByte)
	return &mailConf, nil
}

func (m *MailData) WriteConfig(configFile string) error {
	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	m.Password = base64.StdEncoding.EncodeToString([]byte(m.Password))
	log.Infof("writing config file %s", configFile)
	w := bufio.NewWriter(f)
	e := toml.NewEncoder(w)
	err = e.Encode(m)
	return err
}
