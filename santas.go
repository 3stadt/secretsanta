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

	// filter out mailcontent. TODO: find more clean solution
	// filter out smtp-data. TODO: find more clean solution
	cleanKeys := [][]byte{}
	for _, key := range keys {
		if string(key) != "mailContent" && string(key) != "smtpData" {
			cleanKeys = append(cleanKeys, key)
		}
	}

	res := slowpoke.Gets(c.santaDb, cleanKeys)
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

func pair(s []santa, seed *int64) (map[santa]santa, *int64) {
	if seed == nil {
		now := time.Now().UnixNano()
		seed = &now
	}
	s = removeDuplicates(s)
	rand.Seed(*seed)
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })

	pairs := make(map[santa]santa)

	maxIndex := len(s) - 1
	for i := 0; i <= maxIndex; i++ {
		if i == maxIndex {
			pairs[s[i]] = s[0]
			continue
		}
		pairs[s[i]] = s[i+1]
	}

	return pairs, seed
}

func removeDuplicates(elements []santa) []santa {
	// Use map to record duplicates as we find them.
	santaMap := make(map[santa]struct{})
	result := []santa{}

	for _, v := range elements {
		if v.Name == "" || v.Mail == "" {
			continue
		}
		santaMap[v] = struct{}{}
	}

	for v := range santaMap {
		result = append(result, v)
	}

	return result
}
