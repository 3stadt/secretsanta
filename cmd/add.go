package cmd

import (
	"errors"
	"github.com/3stadt/secretsanta/santa"
	"github.com/mozillazg/go-slugify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"strings"
)

var importCsv string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds new people/secret santas to the database",
	Long:  `The required argument specifies which list should be used`,
	Args:  cobra.ExactArgs(1),
	Run:   Add,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&importCsv, "input", "i", "", "Specifies a csv file to import data from")
	addCmd.MarkFlagRequired("input")
}

func Add(cmd *cobra.Command, args []string) {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(args[0]))
		if b == nil {
			return errors.New("list does not exist yet")
		}
		err := b.ForEach(func(k, v []byte) error {
			log.Infof("A %s is %s.\n", k, v)
			return nil
		})
		return err
	})
	if err != nil {
		log.Warn(err)
	}

	participants, err := santa.Import(importCsv)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%#v\n\n", participants)

	err = db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(slugify.Slugify(args[0])))
		if err != nil {
			log.Fatal(err)
		}
		for _, participant := range *participants {
			key := strings.ToLower(participant.Email)
			err := b.Put([]byte(key), []byte(participant.Name))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
