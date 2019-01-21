package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/recoilme/slowpoke"
	"math/rand"
	"net/http"
	"time"
)

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
