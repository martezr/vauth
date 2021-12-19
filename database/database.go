package database

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

func StartDB(dbdir string) (database *bolt.DB) {
	dbpath := dbdir + "/vauth.db"
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("VirtualMachines"))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func AddDBRecord(db *bolt.DB, key string, data string) {
	// store some data
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("VirtualMachines"))
		if err != nil {
			return err
		}

		log.Printf("Record Added: %s", data)
		err = bucket.Put([]byte(key), []byte(data))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func ViewDBRecord(db *bolt.DB, key string) (data string) {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("VirtualMachines"))
		v := b.Get([]byte(key))
		data = string(v)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Record Read: %s", data)
	return data
}

func ListDBRecords(db *bolt.DB) {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("VirtualMachines"))
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
