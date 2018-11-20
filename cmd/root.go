package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"os"
)

var db *bolt.DB

var rootCmd = &cobra.Command{
	Use:   "sesa",
	Short: "SecretSanta is a cli tool for matching secret santas",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		db, err = bolt.Open("./data.bdb", 0666, nil)
		if err != nil {
			log.Fatal(err)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.Close()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
