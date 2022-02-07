package database

import (
	"log"

	"github.com/martezr/vauth/utils"
	bolt "go.etcd.io/bbolt"
)

// StartDB instantiates the database
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

// AddDBRecord adds a database record
func AddDBRecord(db *bolt.DB, key string, data string) {
	// store some data
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("VirtualMachines"))
		if err != nil {
			return err
		}

		//log.Printf("Record Added: %s", data)
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

// ViewDBRecord gets a single database record
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
	//log.Printf("Record Read: %s", data)
	return data
}

// ListDBRecords gets all the database records
func ListDBRecords(db *bolt.DB) (vms []utils.VMRecord) {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("VirtualMachines"))
		b.ForEach(func(k, v []byte) error {
			var vmdata utils.VMRecord
			vmdata.Name = string(k)
			vmdata.LatestEventId = string(v)
			vms = append(vms, vmdata)
			return nil
		})

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return vms
}
