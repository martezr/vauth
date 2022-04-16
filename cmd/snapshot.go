package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

func init() {
	rootCmd.AddCommand(snapshotCmd)
}

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Create a database snapshot",
	Long:  `Show this help output, or the help for a specified subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
		Snapshot()
		log.Println("snapshot created")
	},
}

func Snapshot() {
	dbpath := "vauth.db"
	db, dberr := bolt.Open(dbpath, 0600, nil)
	if dberr != nil {
		log.Fatal(dberr)
	}
	err := db.View(func(tx *bolt.Tx) error {
		f, err := os.Create("mydb.snap")
		if err != nil {
			log.Println(err)
		}

		defer f.Close()
		_, err = tx.WriteTo(f)
		return err
	})
	if err != nil {
		log.Println(err)
	}
}

func BackupHandleFunc(w http.ResponseWriter, req *http.Request) {
	err := db.View(func(tx *bolt.Tx) error {
		f, err := os.Create("mydb.snap")
		if err != nil {
			log.Println(err)
		}

		defer f.Close()
		//w.Header().Set("Content-Type", "application/octet-stream")
		//w.Header().Set("Content-Disposition", `attachment; filename="my.snap"`)
		//w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err = tx.WriteTo(f)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
